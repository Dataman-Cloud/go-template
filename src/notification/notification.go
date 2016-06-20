package notification

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Dataman-Cloud/go-template/src/config"
	log "github.com/Sirupsen/logrus"
)

var engine_ *NotificationEngine
var once sync.Once

const (
	MessageStaleTime   = time.Duration(time.Minute * 10) // message not sent after this value marked as stale
	RetryAfer          = time.Duration(time.Second * 2)  // initial delay if message not sent successfully
	RetryBackoffFactor = 2                               // backoff factor for a retry
	MessageGCTime      = time.Minute * 10                // 每隔一段时间清理下unsend的消息 interval
	MessageDeleteTime  = time.Hour * 24
)

const (
	SendingChanSize     = 1 << 10
	SinkDumpingChanSize = 1 << 10
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
	StopChan          chan struct{}
}

type Message struct {
	Id           uint64    `json:"id" db:"id"`
	Type         string    `json:"message_type" db:"message_type"`
	ResourceId   uint64    `json:"resource_id" db:"resource_id"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	Time         time.Time `json:"time" db:"time"`
	SinkName     string    `json:"sink_name" db:"sink_name"`
	DumpBegin    time.Time `json:"dump_time" db:"dump_time"`
	Persisted    bool
}

type NotificationEngine struct {
	Running     bool
	Sinks       []*Sink
	SendingChan chan *Message
	StopChan    chan struct{}
}

func NewEngine() *NotificationEngine {
	if engine_ != nil {
		return engine_
	}

	once.Do(func() {
		engine_ = &NotificationEngine{
			Running:     false,
			SendingChan: make(chan *Message, SendingChanSize),
			StopChan:    make(chan struct{}),
			Sinks:       make([]*Sink, 0),
		}
	})

	return engine_
}

func (engine *NotificationEngine) LoadSinks() error {

	sinks := strings.Split(config.GetConfig().Notification, "|")

	for _, value := range sinks {
		urlSink, err := url.Parse(strings.TrimSpace(value))
		if err != nil {
			log.Fatal("sink config error")
		}

		v := urlSink.Query()

		sink := &Sink{Name: v.Get("name"),
			Url:               urlSink.Scheme + "://" + urlSink.Host + urlSink.Path,
			Mode:              v.Get("mode"),
			NotificationTypes: v.Get("notification_types"),
			DumpChan:          make(chan *Message, SinkDumpingChanSize),
			StopChan:          make(chan struct{}),
		}
		engine_.Sinks = append(engine_.Sinks, sink)
	}

	return nil
}
func (engine *NotificationEngine) Stop() {

	close(engine.StopChan)

}
func (engine *NotificationEngine) Start() error {
	if engine.Running {
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
		case msg := <-engine.SendingChan:
			log.Infof("msg sink name %s", msg.SinkName)
			for _, sink := range engine.Sinks {
				if strings.Contains(sink.NotificationTypes, msg.Type) {

					// 新来的消息 或者从数据库里读出属于该sink的
					if msg.SinkName == "" {
						copyMsg := CopyMessage(msg)
						copyMsg.SinkName = sink.Name
						sink.DumpChan <- copyMsg
					} else if msg.SinkName == sink.Name {
						log.Infof("sink name :%s, msg sink_name %s", sink.Name, msg.SinkName)
						sink.DumpChan <- msg
					}

				}
			}
		case <-time.Tick(MessageGCTime):
			go engine.PeriodicallyrMessageGC()

		case <-engine.StopChan:
			for _, sink := range engine.Sinks {
				close(sink.StopChan)
			}
			log.Info("Notification stoped")
			return nil
		}
	}
	return nil
}

// TODO
func (engine *NotificationEngine) HandleStaleMessages() {
	//获取10分钟以内的消息发送
	msgs := LoadMessages(time.Now().Add(MessageStaleTime*-10), time.Now())
	for i, _ := range msgs {
		//	log.Infof("load sinkname %s", msg.SinkName)
		engine.Write(&msgs[i])
	}
}

func (engine *NotificationEngine) PeriodicallyrMessageGC() {
	CleanOutdateMessagesBefore(time.Now().Add(MessageDeleteTime * -1))

	//获取10分钟以外的消息发送
	msgs := LoadMessages(time.Now().Add(MessageDeleteTime*-1), time.Now().Add(MessageStaleTime*-1))
	for i, _ := range msgs {
		engine.Write(&msgs[i])
	}

}

func (engine *NotificationEngine) Write(msg *Message) {
	go func(msg *Message) {
		log.Infof("write %s", msg.SinkName)
		engine.SendingChan <- msg
	}(msg)
}

func (sink *Sink) StartDump() {
	for {
		select {
		case msg := <-sink.DumpChan:
			if !msg.Persisted {
				msg.DumpBegin = time.Now()
			}
			if sink.Mode == MODE_STRICT {
				go sink.StrictWrite(msg)
			} else {
				go sink.BestDeliverWrite(msg)
			}
		case <-sink.StopChan:
			log.Infof("Sink %s stoped", sink.Name)
			return
		}
	}
}

func (sink *Sink) BestDeliverWrite(msg *Message) {
	err := sink.HttpPost(msg)
	if err != nil {
		log.Error("dump message failed", msg.Id)
	}

	return
}

func (sink *Sink) StrictWrite(msg *Message) {
	err := sink.HttpPost(msg)
	if err != nil {
		log.Error("dump message failed, retry after ", RetryAfer)
		msg.SinkName = sink.Name
		if !msg.Persisted {
			msg.Persist()
		}
		sink.StrictWriteRetry(msg, RetryAfer)
	} else {

		log.Infof("dump message successfully %v", *msg)
		if msg.Persisted {
			msg.Remove()
		}
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
		log.Errorln("dump message failed, retry after", retryAfter*RetryBackoffFactor)
		sink.StrictWriteRetry(msg, retryAfter*RetryBackoffFactor)
	} else {
		msg.Remove() // send success remove mssage & stop goroutine
	}
}

func (sink *Sink) HttpPost(msg *Message) error {

	body, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Marshal msg. %s", err.Error())
		return err
	}

	request, err := http.NewRequest("POST", sink.Url, strings.NewReader(string(body)))

	if err != nil {
		log.Errorf("NewRequest failed. %s", err.Error())
		return err

	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("Post request failed. %s", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Response failed. code %d", resp.StatusCode)
		err = errors.New("Response failed")
		return err
	}

	return nil
}
