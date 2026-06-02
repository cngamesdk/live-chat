package v1

import (
	"github.com/cngamesdk/live-chat/api/model"
	"github.com/cngamesdk/live-chat/api/ws"
	"github.com/gin-gonic/gin"
)

type SessionApi struct{}

type CloseSessionRequest struct {
	SessionID int64  `json:"session_id" binding:"required"`
	Reason    string `json:"reason"`
}

// CloseSession godoc
// @Summary Close a chat session (called by admin backend)
// @Tags Admin-Session
// @Param body body CloseSessionRequest true "close params"
// @Success 200 {object} model.Response
// @Router /api/v1/admin/close-session [post]
func (a *SessionApi) CloseSession(c *gin.Context) {
	var req CloseSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailWithMessage("invalid params: "+err.Error(), c)
		return
	}

	// Close session in database
	err := ServiceGroupApp.ChatService.CloseSession(c.Request.Context(), req.SessionID)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	// Notify via WebSocket
	reason := req.Reason
	if reason == "" {
		reason = "会话已结束，感谢您的咨询"
	}
	ws.GlobalHub.SessionClosed(req.SessionID)
	ws.GlobalHub.SendSystemMessage(req.SessionID, reason)

	model.OkWithMessage("success", c)
}
