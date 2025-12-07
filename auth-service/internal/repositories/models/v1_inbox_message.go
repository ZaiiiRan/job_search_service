package models

import (
	"encoding/json"
	"time"
)

type V1InboxMessage struct {
	Id          int64           `db:"id" json:"id"`
	MessageType string          `db:"message_type" json:"message_type"`
	Payload     json.RawMessage `db:"payload" json:"payload"`
	Status      string          `db:"status" json:"status"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}

func (m V1InboxMessage) IsNull() bool { return false }
func (m V1InboxMessage) Index(i int) any {
	switch i {
	case 0:
		return m.Id
	case 1:
		return m.MessageType
	case 2:
		return m.Payload
	case 3:
		return m.Status
	case 4:
		return m.CreatedAt
	case 5:
		return m.UpdatedAt
	default:
		return nil
	}
}
