package notification

import (
	"errors"
	"fmt"
	"time"

	"github.com/Dataman-Cloud/go-template/src/db"
	log "github.com/Sirupsen/logrus"
)

// create a message
func NewMessage() *Message {
	message := &Message{}
	return message
}

func CopyMessage(msg *Message) *Message {
	message := &Message{}

	message.Id = msg.Id
	message.Type = msg.Type
	message.ResourceId = msg.ResourceId
	message.ResourceType = msg.ResourceType
	return message
}

func CleanTooOldMessage(duration time.Duration) error {

	db := db.DB()

	dump_time := time.Now().Add(duration * (-1))

	_, err := db.NamedExec(
		"delete from message where dump_time < :dump_time",
		map[string]interface{}{
			"dump_time": dump_time,
		},
	)

	if err != nil {
		err = errors.New("Remove messgae error: " + err.Error())
		return err
	}

	return nil
}

func LoadMessagesAfter(duration time.Duration) []Message {
	msgs := []Message{}
	db := db.DB()

	sql := fmt.Sprintf(`select * from message where dump_time < %d`, time.Now().Add(duration*(-1)))

	err := db.Select(&msgs, sql)
	if err != nil {

		log.Errorln(" Query from db error: " + err.Error())

	}

	return msgs
}

// load unsent messages from storage
func LoadMessagesBefore(duration time.Duration) []Message {
	// initilize all message marked as persisted

	msgs := []Message{}
	db := db.DB()

	sql := fmt.Sprintf(`select * from message where dump_time > %d`, time.Now().Add(duration*(-1)))

	err := db.Select(&msgs, sql)
	if err != nil {
		log.Errorln(" Query from db error: " + err.Error())
	}

	return msgs
}

// persist a message into storage
func (message *Message) Persist() error {

	message.Persisted = true

	db := db.DB()
	sql := `insert into message(id, type, resource_id, resource_type,sink_name,dump_time) 
	values(:id, :type, :resource_id, :resource_type,:sink_name,:dump_time)`
	_, err := db.NamedExec(sql, message)
	if err != nil {
		err = errors.New("Insert messgae error: " + err.Error())
		return err
	}

	return nil

}

// remove a message permanantly from storage
func (message *Message) Remove() error {
	db := db.DB()
	_, err := db.NamedExec(
		"delete from message where id = :id and sink_name = :sink_name",
		map[string]interface{}{
			"id":        message.Id,
			"sink_name": message.SinkName,
		},
	)

	if err != nil {
		err = errors.New("Remove messgae error: " + err.Error())
		return err
	}

	return nil
}
