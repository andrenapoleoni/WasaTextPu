package api

import (
	"encoding/json"
	"myWasatext/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// check authorization
	profileUserID, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid user id")
		return
	}

	userID := ctx.UserID

	if profileUserID != userID {
		Forbidden(w, err, ctx, "Forbidden 1")
		return
	}

	// take chat id from endpoint
	ConversationID, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid conversation id")
		return
	}
	// check if conversation exists

	conversation, err := rt.db.GetConversation(ConversationID)
	if err != nil {
		http.Error(w, "Not Found "+err.Error(), http.StatusNotFound)
		return
	}

	// check if user is in conversation
	if conversation.GroupID != 0 {
		exist, err := rt.db.ExistUserInGroup(userID, conversation.GroupID)
		if err != nil {
			InternalServerError(w, err, "Failed to check if user is in group", ctx)
			return
		}
		if !exist {
			Forbidden(w, err, ctx, "User is not in group")
			return
		}
	} else {
		exist, err := rt.db.ExistUserInConv(userID, conversation.ConversationID)
		if err != nil {
			InternalServerError(w, err, "Failed to check if user is in conversation", ctx)
			return
		}
		if !exist {
			Forbidden(w, err, ctx, "User is not in conversation")
			return
		}
	}

	// take message id from endpoint
	MessageID, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid message id")
		return
	}
	// check if message exists
	message, err := rt.db.GetMessage(conversation.ConversationID, MessageID)
	if err != nil {
		http.Error(w, "Not Found "+err.Error(), http.StatusNotFound)
		return
	}
	message.UserID = userID
	// take dest id from query (conv id)
	dest := r.URL.Query().Get("dest")
	if dest == "" {
		BadRequest(w, err, ctx, "Invalid dest id")
		return
	}

	//check if dest is a username or a int
	destID, err := strconv.Atoi(dest)
	if err != nil {
		//if is a username get id
		userRec, err := rt.db.GetUserByName(dest)
		if err != nil {
			InternalServerError(w, err, "Failed to get user by username", ctx)
			return
		}
		// check if conversation between users already exists
		exist, err := rt.db.ExistConversation(userID, userRec.UserID)
		if err != nil {
			InternalServerError(w, err, "Failed to check if conversation exists", ctx)
			return
		}
		if exist {
			//get conversation id
			convReceiverDB, err := rt.db.GetConversationIDfrom2Users(userID, userRec.UserID)
			if err != nil {
				InternalServerError(w, err, "Failed to get conversation id", ctx)
				return
			}

			message.ConversationID = convReceiverDB.ConversationID
			message.Forwarded = true
			msg, err := rt.db.CreateMessage(message)
			if err != nil {
				InternalServerError(w, err, "Failed to create message", ctx)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			if err := json.NewEncoder(w).Encode(msg); err != nil {
				InternalServerError(w, err, "Failed to encode response", ctx)
				return
			}
		} else {
			//create conversation
			//if not exist create conversation and send message
			var convReceiver Conversation
			convReceiver.GroupID = 0
			convReceiver, err = rt.CreateConversationDB(convReceiver)
			if err != nil {
				InternalServerError(w, err, "Failed to create conversation", ctx)
				return
			}
			//add members to conversation
			err = rt.db.AddMemberPrivate(convReceiver.ConversationID, userID)
			if err != nil {
				InternalServerError(w, err, "Failed to add user to conversation", ctx)
				return
			}
			err = rt.db.AddMemberPrivate(convReceiver.ConversationID, userRec.UserID)
			if err != nil {
				InternalServerError(w, err, "Failed to add user to conversation", ctx)
				return
			}
			//send message
			message.ConversationID = convReceiver.ConversationID
			message.Forwarded = true
			msg, err := rt.db.CreateMessage(message)
			if err != nil {
				InternalServerError(w, err, "Failed to create message", ctx)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			if err := json.NewEncoder(w).Encode(msg); err != nil {
				InternalServerError(w, err, "Failed to encode response", ctx)
				return
			}

		}

	} else {

		//is a group
		//take convid from groupId

		convGrDB, err := rt.db.GetConversationIDfromGroup(destID)
		if err != nil {
			InternalServerError(w, err, "Failed to get conversation", ctx)
			return
		}
		message.ConversationID = convGrDB.ConversationID
		message.Forwarded = true
		msg, err := rt.db.CreateMessage(message)
		if err != nil {
			InternalServerError(w, err, "Failed to create message", ctx)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			InternalServerError(w, err, "Failed to encode response", ctx)
			return
		}

	}
}
