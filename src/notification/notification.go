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
	events            chan Message
}

type Message struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"`
	ResourceId   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Time         time.Time `json:"time"`
	SinkName     string
	RetryTime    uint64
}

type NotificationEngine struct {
	Runing bool
	Sinks  []*Sink
	events chan Message
	unsend chan struct{}
}

func NewEngine() *NotificationEngine {
	if engine_ != nil {
		return engine_
	}

	once.Do(func() {
		engine_ = &NotificationEngine{
			Runing: false,
			events: make(chan Message, 999),
			Sinks:  make([]*Sink, 0),
		}
		// TODO
		// 1 parseSinks
		sink1 := &Sink{Name: "Name1",
			Url:               "http://localhost:8080/foo/bar",
			Mode:              "strict",
			NotificationTypes: MESSAGE_TYPE_APP_DELETION,
			events:            make(chan Message, 999)}
		sink2 := &Sink{Name: "Name2",
			Url:               "http://localhost:8089/foo/bar",
			Mode:              "strict",
			NotificationTypes: MESSAGE_TYPE_APP_DELETION,
			events:            make(chan Message, 999)}

		engine_.Sinks = append(engine_.Sinks, sink1)
		engine_.Sinks = append(engine_.Sinks, sink2)

		for _, sink := range engine_.Sinks {
			go sink.run()
		}

		go engine_.Start()

	})

	return engine_
}

func (engine *NotificationEngine) Start() error {
	if engine.Runing {
		log.Infoln("NotificationEngine already start")
		return nil
	}

	//启动时发送数据库里的消息
	go engine.RetryUnSend()

	for {
		select {
		case msg := <-engine.events:
			for _, sink := range engine.Sinks {

				if strings.Contains(sink.NotificationTypes, msg.Type) {
					continue
				}
				// 为空为第一次发的消息  sinkname有值是retry
				if msg.SinkName == "" || msg.SinkName == sink.Name {
					sink.events <- msg
				}
			}
		case <-time.After(time.Second * 30):
			go engine.RetryUnSend()
		}
	}

	return nil
}

func (engine *NotificationEngine) RetryUnSend() {

	msgs := LoadMessages()

	for _, msg := range msgs {
		msg.Remove()
		engine.Write(msg)
	}

}

func (engine *NotificationEngine) Write(msg Message) {
	go func(msg Message) {
		engine.events <- msg
	}(msg)
}

func (sink *Sink) Write(msg Message) error {

	var err error
	if sink.Mode == MODE_STRICT {
		err = sink.strictWrite(msg)
	} else {
		err = sink.delivery(msg)
	}
	return err
}

func (sink *Sink) run() {
	for {
		select {
		case msg := <-sink.events:
			if err := sink.Write(msg); err != nil {
				log.Errorf("Failed to send msg to %s. %s", sink.Name, err.Error())
			}

		}
	}
}

func (sink *Sink) delivery(msg Message) error {
	//发送到Sink url
	return nil
}

func (sink *Sink) strictWrite(msg Message) error {

	err := sink.delivery(msg)

	if err != nil {

		//		if msg.RetryTime > 5 {
		//			log.Errorln("Delete msg after 5 retry")
		//			return err
		//		}

		msg.SinkName = sink.Name
		msg.Time = time.Now()
		msg.RetryTime++
		msg.Persist()
	}

	return err
}
