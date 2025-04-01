package api

import (
	"encoding/json"
	"myWasatext/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// get conversationID from endpoint
	conversationID, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// check if conversation exists
	exist, err := rt.db.ExistConversationByID(conversationID)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error 1", ctx)
		return
	}
	if !exist {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var membergroup []User
	// get conversation to check if groupId !=0
	conversation, err := rt.db.GetConversation(conversationID)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error 2", ctx)
		return
	}
	if conversation.GroupID != 0 {
		allusers, err := rt.db.GetMemberGroup(conversation.GroupID)

		if err != nil {
			InternalServerError(w, err, "can't get the list of users", ctx)
			return
		}
		for _, u := range allusers {
			usDB, err := rt.db.GetUserByID(u)
			if err != nil {
				InternalServerError(w, err, "can't get the user", ctx)
				return
			}
			var us User
			err = us.FromDatabase(usDB)
			if err != nil {
				InternalServerError(w, err, "can't get the user", ctx)
				return
			}

			membergroup = append(membergroup, us)

		}

	}

	// get messages of conversation
	messageDB, err := rt.db.GetAllMessage(conversationID)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error 3", ctx)
		return
	}
	type MexR struct {
		MessageID      int    `json:"messageID"`
		ConversationID int    `json:"conversationID"`
		UserID         int    `json:"userID"`
		MessageTXT     string `json:"txt"`
		Timestamp      string `json:"time"`
		Forwarded      bool   `json:"forwarded"`
		Photo          string `json:"photo"`
		Linkmessage    int    `json:"linkmessage"`
		Checkmark      bool   `json:"checkmark"`
	}
	type CommentR struct {
		CommentTXT string `json:"commentTXT"`
		Username   string `json:"username"`
		CommentId  int    `json:"commentId"`
	}

	type MessageResponse struct {
		MessageR    MexR       `json:"message"`
		UserR       User       `json:"user"`
		CommentResp []CommentR `json:"comment"`
	}

	type MessageListResponse struct {
		Messages   []MessageResponse `json:"messages"`
		MemberList []User            `json:"memberlist"`
	}

	response := make([]MessageResponse, len(messageDB))
	for i, message := range messageDB {
		userDB, err := rt.db.GetUserByID(message.UserID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error 4", ctx)
			return
		}
		var user User
		err = user.FromDatabase(userDB)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error 5", ctx)
			return
		}
		var rsp MessageResponse
		var msg Message
		//fmt.Println("userID", userID, "message.UserID", message.UserID)
		if message.UserID != userID {
			err = rt.db.DeleteCheckmark(message.MessageID, message.ConversationID, userID)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error 6", ctx)
				return
			}

		}
		exist, err := rt.db.ExistCheckrow(message.MessageID, message.ConversationID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error 6", ctx)
			return
		}

		if !exist {
			err := rt.db.UpdateMessage(message.MessageID, message.ConversationID)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error 7", ctx)
				return
			}

		}

		err = msg.FromDatabase(message)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error 6", ctx)
			return
		}
		var mex MexR
		mex.MessageID = msg.MessageID
		mex.ConversationID = msg.ConversationID
		mex.UserID = msg.UserID
		mex.MessageTXT = msg.MessageTXT
		mex.Forwarded = msg.Forwarded
		mex.Photo = msg.Photo
		mex.Linkmessage = msg.Linkmessage
		mex.Timestamp = message.Timestamp.Format("15:04")
		mex.Checkmark = msg.Checkmark

		var Comments []CommentR
		commentsDB, err := rt.db.GetComments(message.MessageID, message.ConversationID)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error 7", ctx)
			return
		}

		for _, comment := range commentsDB {
			var CommentR CommentR
			CommentR.CommentTXT = comment.CommentTXT
			CommentR.CommentId = comment.CommentID

			userDB, err := rt.db.GetUserByID(comment.UserID)
			if err != nil {
				InternalServerError(w, err, "Internal Server Error 8", ctx)
				return
			}
			CommentR.Username = userDB.Username

			Comments = append(Comments, CommentR)
		}
		rsp.CommentResp = Comments
		rsp.MessageR = mex
		rsp.UserR = user

		response[i] = rsp
	}

	finalResponse := MessageListResponse{
		Messages:   response,
		MemberList: membergroup,
	}
	// return
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(finalResponse)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error 9", ctx)
		return
	}

}
