package api

import (
	"encoding/json"
	"myWasatext/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// check authorization
	profileUserID, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	userID := ctx.UserID

	// check authorization
	if profileUserID != userID {
		Forbidden(w, err, ctx, "Forbidden")
		return
	}
	// get conversations
	conversations, err := rt.db.GetListConversations(userID)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error1", ctx)
		return
	}
	type MexR struct {
		MessageTxt string `json:"mexTXT"`
		Username   string `json:"mexUser"`
		Timestamp  string `json:"mexTime"`
	}

	// convert from conv id to username or groupname
	type Response struct {
		Conversation Conversation `json:"conversation"`
		User         User         `json:"user"`
		Group        Group        `json:"group"`
		GroupUsers   []User       `json:"groupUsers"`
		LastMessage  MexR         `json:"lastmessage"`
	}
	// Response
	response := make([]Response, len(conversations))

	for i, conversations := range conversations {
		if conversations.GroupID != 0 {
			groupDB, err := rt.db.GetGroupByID(conversations.GroupID)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error2", ctx)
				return
			}
			var group Group
			err = group.FromDatabase(groupDB)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error3", ctx)
				return
			}
			// take all users in group
			groupUsers, err := rt.db.GetUsersInGroup(conversations.GroupID)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error4", ctx)
				return
			}
			var groupUsersResponse []User
			for i := 0; i < len(groupUsers); i++ {
				userD, err := rt.db.GetUserByID(groupUsers[i])
				if err != nil {
					InternalServerError(w, err, "Internal Server Error5", ctx)
					return
				}
				var user User
				err = user.FromDatabase(userD)
				if err != nil {
					InternalServerError(w, err, "Internal Server Error6", ctx)
					return
				}

				groupUsersResponse = append(groupUsersResponse, user)
			}

			response[i].Group = group
			response[i].GroupUsers = groupUsersResponse
		} else {
			userDB, err := rt.db.GetUserInConversationPrivate(conversations.ConversationID, userID)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error4", ctx)
				return
			}
			var user User
			err = user.FromDatabase(userDB)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error5", ctx)
				return
			}

			response[i].User = user
		}
		var conversation Conversation
		err = conversation.FromDatabase(conversations)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error6", ctx)
			return
		}

		response[i].Conversation = conversation

		message, err := rt.db.GetLastMessage(conversation.ConversationID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error 7", ctx)
			return
		}
		var mex MexR
		mex.MessageTxt = message.MessageTXT
		mex.Timestamp = message.Timestamp.Format("2006-01-02 - 15:04")
		user, err := rt.db.GetUserByID(message.UserID)
		if err != nil && err.Error() != "sql: no rows in result set" {
			InternalServerError(w, err, "Internal Server Error 8", ctx)
			return
		}
		mex.Username = user.Username

		response[i].LastMessage = mex

	}

	// return conversations
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error9", ctx)
		return
	}

}
