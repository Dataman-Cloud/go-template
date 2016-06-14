package notification

import (
	"fmt"
	"sync"

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
}

type Message struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"` //
	ResourceId   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Time         time.Time `json:"time"`
	Sink         *Sink
}

type NotificationEngine struct {
	Runing bool

	Sinks []*Sink
}

func NewEngine() *NotificationEngine {
	if engine_ {
		return engine_
	}

	once.Do(func() {
		engine_ = &NotificationEngine{
			Runing: false,
		}
		// TODO
		// 1 parseSinks
		// 2 init channels - send channel, buf channel, retry channels

	})

	return engine_
}
func (engine *NotificationEngine) Start() error {
	if engine.Runing {
		log.Infoln("NotificationEngine already start")
		return nil
	}

	// TODO heavy stuffs go here
	return nil
}
