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

func CleanOutdateMessagesBefore(t time.Time) error {
	db := db.DB()
	_, err := db.NamedExec(
		"DELETE FROM message WHERE dump_time < :dump_time",
		map[string]interface{}{
			"dump_time": t,
		},
	)

	if err != nil {
		err = errors.New("Remove messgae error: " + err.Error())
		return err
	}

	return nil
}

func LoadMessages(from time.Time, to time.Time) []Message {
	msgs := []Message{}
	db := db.DB()

	sql := `SELECT * FROM message WHERE dump_time > '%s' AND dump_time < '%s'`
	sql = fmt.Sprintf(sql, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
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
	sql := `INSERT INTO message(id, message_type, resource_id, resource_type, sink_name, dump_time) 
	VALUES(:id, :type, :resource_id, :resource_type,:sink_name,:dump_time)`
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
		"DELETE FROM message WHERE id = :id and sink_name = :sink_name",
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
