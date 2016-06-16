package notification

import (
	"github.com/Dataman-Cloud/go-template/src/db"
)

// create a message
func NewMessage() *Message {
	message := &Message{}
	return message
}

// load unsent messages from storage
func LoadMessages() []*Message {
	// initilize all message marked as persisted

	//	msgs := make([]*Message, 0)
	//	db := db.DB()
	//	var rows *sqlx.Rows
	//	var err error
	//	sql = `select * from message`

	//	rows, err = db.NamedQuery(sql)
	//	if err != nil {
	//		err := errors.New(" Query from db error: " + err.Error())
	//		return projects, err
	//	}
	//	defer rows.Close()

	//	for rows.Next() {

	//		var id int64
	//		var uid int64
	//		var name string
	//		var pubkey string
	//		var imageName string
	//		var description string
	//		var branch string
	//		var period int64
	//		var repo_uri string
	//		var trigger_type uint8
	//		var active bool
	//		var created time.Time
	//		var updated time.Time
	//		if err = rows.Scan(&id, &uid, &period, &name, &imageName, &description, &pubkey, &branch, &repo_uri, &trigger_type, &active, &created, &updated); err != nil {
	//			err := errors.New("[ListProject] scan from rows errors: " + err.Error())
	//			return projects, err
	//		}

	//		msg := Message{
	//			Id:          id,
	//			Uid:         uid,
	//			Period:      period,
	//			Name:        name,
	//			ImageName:   imageName,
	//			Description: description,
	//			Pubkey:      pubkey,
	//			Branch:      branch,
	//			RepoUri:     repo_uri,
	//			TriggerType: trigger_type,
	//			Active:      active,
	//			Created:     created.Unix(),
	//			Updated:     updated.Unix(),
	//		}
	//		msgs = append(msgs, msg)
	//	}
	//	return nil
}

// persist a message into storage
func (message *Message) Persist() error {
	message.Persisted = true

	db := db.DB()
	sql := `insert into message(id, type, resource_id, resource_type) values(:id, :type, :resource_id, :resource_type)`
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
		"delete from message where id = :id",
		map[string]interface{}{
			"id": message.Id,
		},
	)

	if err != nil {
		err = errors.New("Remove messgae error: " + err.Error())
		return err
	}

	return nil
}
