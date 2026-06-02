package v1

import (
	"time"

	"github.com/cngamesdk/live-chat/api/model"
	"github.com/cngamesdk/live-chat/api/ws"
	"github.com/gin-gonic/gin"
)

type AdminApi struct{}

type AdminAgentReplyRequest struct {
	SessionID int64  `json:"session_id" binding:"required"`
	AgentID   int64  `json:"agent_id"`
	AgentName string `json:"agent_name"`
	Content   string `json:"content" binding:"required"`
	MsgType   string `json:"msg_type"`
}

// AgentReply godoc
// @Summary Agent reply to user (called by admin backend)
// @Tags Admin-Chat
// @Param body body AdminAgentReplyRequest true "reply params"
// @Success 200 {object} model.Response
// @Router /api/v1/admin/agent-reply [post]
func (a *AdminApi) AgentReply(c *gin.Context) {
	var req AdminAgentReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailWithMessage("invalid params: "+err.Error(), c)
		return
	}

	if req.MsgType == "" {
		req.MsgType = "text"
	}
	if req.AgentName == "" {
		req.AgentName = "客服"
	}

	// Save message to database
	msg, err := ServiceGroupApp.ChatService.SaveAgentMessage(c.Request.Context(), req.SessionID, req.AgentID, req.AgentName, req.Content, req.MsgType)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	// Broadcast via WebSocket
	ws.GlobalHub.SendToSession(req.SessionID, &ws.Message{
		Type:       "message",
		Content:    msg.Content,
		MsgType:    msg.MsgType,
		SenderType: msg.SenderType,
		SenderID:   msg.SenderID,
		SenderName: msg.SenderName,
		MessageID:  msg.ID,
		IsRead:     msg.IsRead,
		Timestamp:  time.Now().UnixMilli(),
	})

	model.OkWithData(gin.H{"message_id": msg.ID}, c)
}
