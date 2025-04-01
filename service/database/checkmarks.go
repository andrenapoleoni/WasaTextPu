package database

import "database/sql"

var query_ADDLINECHECKMARK = `INSERT INTO Checkmark (messageID, conversationID, userID)
						VALUES  (?, ?, ?)`

var query_DELETECHECKMARK = `DELETE FROM Checkmark WHERE messageID = ? AND conversationID = ? AND userID = ?`

var query_EXISTROW = `SELECT EXISTS(SELECT 1 FROM Checkmark WHERE messageID = ? AND conversationID = ?)`

var query_UPDATECHEK = `UPDATE Message SET Checkmark = true WHERE messageID = ? AND conversationID = ?`

func (db *appdbimpl) AddCheckmark(messageID int, conversationID int, userID int) error {
	_, err := db.c.Exec(query_ADDLINECHECKMARK, messageID, conversationID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (db *appdbimpl) DeleteCheckmark(messageID int, conversationID int, userID int) error {
	_, err := db.c.Exec(query_DELETECHECKMARK, messageID, conversationID, userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}
func (db *appdbimpl) ExistCheckrow(messageID int, conversationID int) (bool, error) {
	var exist bool
	err := db.c.QueryRow(query_EXISTROW, messageID, conversationID).Scan(&exist)
	if err != nil && err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (db *appdbimpl) UpdateMessage(messageID int, conversationID int) error {
	_, err := db.c.Exec(query_UPDATECHEK, messageID, conversationID)
	if err != nil {
		return err
	}
	return nil

}
