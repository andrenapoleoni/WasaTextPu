package database

import (
	"database/sql"
	"errors"
)

var query_GETMESSAGE = `SELECT message, userID,forwarded,photo FROM Message WHERE messageID = ? AND conversationID = ?;`
var query_GETLASTMESSAGE = `SELECT message, userID,TIMESTAMP FROM Message WHERE conversationID = ? ORDER BY messageID DESC LIMIT 1;`

func (db *appdbimpl) GetMessage(conversationID int, messageID int) (Message, error) {
	// get message from database
	var message Message
	message.MessageID = messageID
	message.ConversationID = conversationID

	err := db.c.QueryRow(query_GETMESSAGE, messageID, conversationID).Scan(&message.MessageTXT, &message.UserID, &message.Forwarded, &message.Photo)
	if err != nil {
		return message, err
	}

	return message, nil

}

func (db *appdbimpl) GetAllMessage(conversationID int) ([]Message, error) {
	// get all message from database
	var messages []Message
	rows, err := db.c.Query("SELECT messageID, message, userID,forwarded,TIMESTAMP,photo,linkmessage,checkmark FROM Message WHERE conversationID = ?;", conversationID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		err := rows.Scan(&message.MessageID, &message.MessageTXT, &message.UserID, &message.Forwarded, &message.Timestamp, &message.Photo, &message.Linkmessage, &message.Checkmark)
		if err != nil {
			return messages, err
		}
		if rows.Err() != nil {
			return messages, err
		}
		message.ConversationID = conversationID
		messages = append(messages, message)
	}

	return messages, nil

}

func (db *appdbimpl) GetLastMessage(conversationID int) (Message, error) {
	var last Message
	err := db.c.QueryRow(query_GETLASTMESSAGE, conversationID).Scan(&last.MessageTXT, &last.UserID, &last.Timestamp)
	if err != nil && err != sql.ErrNoRows {
		return last, err
	}

	return last, nil
}
