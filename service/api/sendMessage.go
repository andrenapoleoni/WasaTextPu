package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"myWasatext/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreateMessageDB(message Message) (Message, error) {
	// create message in database

	messageDB, err := rt.db.CreateMessage(message.ToDatabase())
	if err != nil {
		return message, err
	}

	// convert message from database
	err = message.FromDatabase(messageDB)
	if err != nil {
		return message, err
	}

	return message, nil

}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// check authorization
	profileUserID, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	userID := ctx.UserID

	if profileUserID != userID {
		Forbidden(w, err, ctx, "auth problem")
		return
	}
	// get conversation id
	conversationID, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}
	// chech if conversation exist
	conversation, err := rt.db.GetConversation(conversationID)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error", ctx)
		return
	}
	// check if user is in conversation
	if conversation.GroupID != 0 {
		exist, err := rt.db.ExistUserInGroup(userID, conversation.GroupID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error", ctx)
			return
		}
		if !exist {
			Forbidden(w, err, ctx, "group not exist")
			return
		}
	} else {
		exist, err := rt.db.ExistUserInConv(userID, conversationID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error", ctx)
			return
		}
		if !exist {
			Forbidden(w, err, ctx, "user not exist")
			return
		}
	}
	var message Message
	// Check if the size of the image is less than 5MB
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		BadRequest(w, err, ctx, "Image too big")
		return
	}
	message.MessageTXT = r.FormValue("text")
	message.Linkmessage, err = strconv.Atoi(r.FormValue("reply"))
	if err != nil && r.FormValue("reply") != "" {
		BadRequest(w, err, ctx, "Bad Request 1")
	}
	file, _, err := r.FormFile("file")

	if message.MessageTXT == "" && file == nil {
		BadRequest(w, err, ctx, "Bad Request 1")
		return
	}
	if err == nil {
		data, err := io.ReadAll(file)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request 2")
			return
		}
		fileType := http.DetectContentType(data)
		if fileType != "image/jpeg" && fileType != "image/gif" {
			BadRequest(w, err, ctx, "Bad Request 3")
			return
		}
		defer func() { err = file.Close() }()
		message.Photo = base64.StdEncoding.EncodeToString(data)
	}

	message.UserID = userID
	message.ConversationID = conversationID
	message.Checkmark = false

	// send message
	message, err = rt.CreateMessageDB(message)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error", ctx)
		return
	}

	// add checkmark table
	if conversation.GroupID != 0 {
		// get users in group
		users, err := rt.db.GetUsersInGroup(conversation.GroupID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error", ctx)
			return
		}
		// add checkmark for each user
		for _, user := range users {
			if user != userID {

				err = rt.db.AddCheckmark(message.MessageID, conversationID, user)
				if err != nil {
					InternalServerError(w, err, "Internal Server Error", ctx)
					return

				}
			}
		}
	} else {
		uzer, err := rt.db.GetUserInConversationPrivate(conversationID, userID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error", ctx)
			return
		}

		err = rt.db.AddCheckmark(message.MessageID, conversationID, uzer.UserID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error", ctx)
			return
		}

	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error", ctx)
		return
	}

}
