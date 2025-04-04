package database

var query_GETCONVERSATION = `SELECT groupID FROM Conversation WHERE conversationID = ?;`

func (db *appdbimpl) GetConversation(conversationID int) (Conversation, error) {
	// get conversation from database
	var conversation Conversation
	conversation.ConversationID = conversationID

	err := db.c.QueryRow(query_GETCONVERSATION, conversationID).Scan(&conversation.GroupID)
	if err != nil {
		return conversation, err
	}

	return conversation, nil

}

func (db *appdbimpl) GetConversationIDfrom2Users(userID1 int, userID2 int) (Conversation, error) {
	// get conversation from database
	var conversation Conversation

	err := db.c.QueryRow(`SELECT c1.conversationID
		FROM MemberPrivate c1
		JOIN MemberPrivate c2 ON c1.conversationID = c2.conversationID
		WHERE c1.userID = ?
		  AND c2.userID = ?;`, userID1, userID2).Scan(&conversation.ConversationID)

	if err != nil {
		return conversation, err
	}

	return conversation, nil

}

func (db *appdbimpl) GetConversationIDfromGroup(groupID int) (Conversation, error) {
	// get conversation from database
	var conversation Conversation

	err := db.c.QueryRow(`SELECT conversationID
		FROM Conversation
		WHERE groupID = ?;`, groupID).Scan(&conversation.ConversationID)

	conversation.GroupID = groupID

	if err != nil {
		return conversation, err
	}

	return conversation, nil

}
