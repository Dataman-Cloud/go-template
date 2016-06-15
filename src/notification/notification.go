package notification

import (
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var engine_ *NotificationEngine
var once sync.Once

const (
	MessageStaleTime   = time.Duration(time.Minute * 10) // message not sent after this value marked as stale
	RetryAfer          = time.Duration(time.Second * 2)  // initial delay if message not sent successfully
	RetryBackoffFactor = 2                               // backoff factor for a retry
)

const (
	SendingChanSize     = 1 << 10
	SinkDumpingChanSIze = 1 << 10
)

const (
	MODE_STRICT        = "strict"
	MODE_BEST_DELIVERY = "best_deliver"
)

const (
	MESSAGE_TYPE_APP_DELETION = "APP_DELETION"
	MESSAGE_TYPE_APP_CREATION = "APP_CREATION"
)

type Sink struct {
	Name              string
	Url               string `yaml:"url"`
	Mode              string `yaml:"mode"`
	NotificationTypes string `yaml:"notification_types"`
	DumpChan          chan *Message
}

type Message struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"`
	ResourceId   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Time         time.Time `json:"time"`
	Sink         *Sink
	Persisted    bool
	DumpBegin    time.Time
}

type NotificationEngine struct {
	Runing      bool
	Sinks       []*Sink
	SendingChan chan *Message
}

func NewEngine() *NotificationEngine {
	if engine_ != nil {
		return engine_
	}

	once.Do(func() {
		engine_ = &NotificationEngine{
			Runing:      false,
			SendingChan: make(chan *Message, SendingChanSize),
			Sinks:       make([]*Sink, 0),
		}
	})

	return engine_
}

func (engine *NotificationEngine) LoadSinks() error {
	sink1 := &Sink{Name: "Name1",
		Url:               "http://localhost:8080/foo/bar",
		Mode:              "strict",
		NotificationTypes: MESSAGE_TYPE_APP_DELETION,
		DumpChan:          make(chan *Message, SinkDumpingChanSIze)}

	sink2 := &Sink{Name: "Name2",
		Url:               "http://localhost:8089/foo/bar",
		Mode:              "strict",
		NotificationTypes: MESSAGE_TYPE_APP_DELETION,
		DumpChan:          make(chan *Message, SinkDumpingChanSIze)}

	engine_.Sinks = append(engine_.Sinks, sink1)
	engine_.Sinks = append(engine_.Sinks, sink2)

	return nil
}

func (engine *NotificationEngine) Start() error {
	if engine.Runing {
		log.Infoln("NotificationEngine already start")
		return nil
	}

	err := engine.LoadSinks()
	if err != nil {
		log.Error("Loading Sink error")
		return err
	}

	for _, sink := range engine_.Sinks {
		go sink.StartDump()
	}

	go engine.HandleStaleMessages()

	for {
		select {
		case msg := <-engine.events:
			for _, sink := range engine.Sinks {
				if strings.Contains(sink.NotificationTypes, msg.Type) {
					msg.Sink = sink
					sink.DumpChan << msg
				}
			}
		}
	}

	return nil
}

// TODO
func (engine *NotificationEngine) HandleStaleMessages() {
	msgs := LoadMessages()

	for _, msg := range msgs {
		msg.Remove()
		engine.Write(msg)
	}

}

func (sink *Sink) Write(msg Message) error {

	var err error
	return err
}

func (sink *Sink) StartDump() {
	for {
		select {
		case msg := <-sink.DumpChan:
			msg.DumpBegin = time.Now()
			if sink.Mode == MODE_STRICT {
				go sink.StrictWrite(msg)
			} else {
				go sink.BestDeliverWrite(msg)
			}
		}
	}
}

func (sink *Sink) BestDeliverWrite(msg *Message) {
	err := sink.HttpPost(msg)
	if err != nil {
		log.Error("dump message failed", msg.Id)
	}

	return nil
}

func (sink *Sink) StrictWrite(msg *Message) {
	err := sink.HttpPost(msg)
	if err != nil {
		log.Error("dump message failed, retry after", msg.Id, retryAfter*RetryBackoffFactor)
		_ := msg.Persist()
		sink.StrictWriteRetry(msg, RetryAfer)
	}
}

func (sink *Sink) StrictWriteRetry(msg *Message, retryAfter time.Duration) {
	if time.Now().Sub(msg.DumpBegin) > MessageStaleTime {
		log.Error("message stable, stop sending", msg.Id)
		return // stop dumping goroutine now
	}

	time.Sleep(retryAfter)
	err := sink.HttpPost(msg)
	if err != nil {
		log.Error("dump message failed, retry after", msg.Id, retryAfter*RetryBackoffFactor)
		sink.StrictWriteRetry(msg, retryAfter*RetryBackoffFactor)
	} else {
		msg.Remove() // send success remove mssage & stop goroutine
	}
}

func (sink *Sink) HttpPost(msg *Message) error {
	return nil
}
