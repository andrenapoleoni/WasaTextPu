package api

import (
	"myWasatext/service/database"
	"time"
)

// Message struct
type Message struct {
	MessageID      int       `json:"messageID"`
	ConversationID int       `json:"conversationID"`
	UserID         int       `json:"userID"`
	MessageTXT     string    `json:"txt"`
	Timestamp      time.Time `json:"time"`
	Forwarded      bool      `json:"forwarded"`
	Photo          string    `json:"photo"`
	Linkmessage    int       `json:"linkmessage"`
	Checkmark      bool      `json:"checkmark"`
}

// from database
func (m *Message) FromDatabase(dbMessage database.Message) error {
	m.MessageID = dbMessage.MessageID
	m.ConversationID = dbMessage.ConversationID
	m.UserID = dbMessage.UserID
	m.MessageTXT = dbMessage.MessageTXT
	m.Timestamp = dbMessage.Timestamp
	m.Forwarded = dbMessage.Forwarded
	m.Photo = dbMessage.Photo
	m.Linkmessage = dbMessage.Linkmessage
	m.Checkmark = dbMessage.Checkmark

	return nil
}

// to database
func (m *Message) ToDatabase() database.Message {
	return database.Message{
		MessageID:      m.MessageID,
		ConversationID: m.ConversationID,
		UserID:         m.UserID,
		MessageTXT:     m.MessageTXT,
		Timestamp:      m.Timestamp,
		Forwarded:      m.Forwarded,
		Photo:          m.Photo,
		Linkmessage:    m.Linkmessage,
		Checkmark:      m.Checkmark,
	}
}
