package v1

import (
	"github.com/cngamesdk/live-chat/api/model"
	"github.com/gin-gonic/gin"
)

type ChatApi struct{}

type InitSessionRequest struct {
	ProductCode string `json:"product_code" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	UserName    string `json:"user_name"`
	UserToken   string `json:"user_token"`
	Source      string `json:"source"`
}

// InitSession godoc
// @Summary Initialize a chat session
// @Tags H5-Chat
// @Param body body InitSessionRequest true "init params"
// @Success 200 {object} model.Response
// @Router /api/v1/chat/init [post]
func (a *ChatApi) InitSession(c *gin.Context) {
	var req InitSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailWithMessage("invalid params: "+err.Error(), c)
		return
	}
	if req.Source == "" {
		req.Source = "h5"
	}
	if req.UserName == "" {
		req.UserName = req.UserID
	}

	session, token, err := ServiceGroupApp.ChatService.InitSession(c.Request.Context(), req.ProductCode, req.UserID, req.UserName, req.Source, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	model.OkWithData(gin.H{
		"session_id": session.ID,
		"token":      token,
		"user_id":    req.UserID,
	}, c)
}

type GetHistoryRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

// GetHistory godoc
// @Summary Get chat history
// @Tags H5-Chat
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {object} model.Response
// @Router /api/v1/chat/history [get]
func (a *ChatApi) GetHistory(c *gin.Context) {
	sessionID := c.GetInt64("session_id")
	var req GetHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 50
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}

	messages, total, err := ServiceGroupApp.ChatService.GetHistory(c.Request.Context(), sessionID, req.Page, req.PageSize)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	model.OkWithPage(messages, total, req.Page, req.PageSize, c)
}

type MarkReadRequest struct {
	MessageIDs []int64 `json:"message_ids" binding:"required"`
}

// MarkMessagesAsRead godoc
// @Summary Mark messages as read
// @Tags H5-Chat
// @Param body body MarkReadRequest true "message ids"
// @Success 200 {object} model.Response
// @Router /api/v1/chat/mark-read [post]
func (a *ChatApi) MarkMessagesAsRead(c *gin.Context) {
	sessionID := c.GetInt64("session_id")
	var req MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.FailWithMessage("invalid params: "+err.Error(), c)
		return
	}

	err := ServiceGroupApp.ChatService.MarkMessagesAsRead(c.Request.Context(), sessionID, req.MessageIDs)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	model.OkWithMessage("success", c)
}

// GetUnreadCount godoc
// @Summary Get unread message count
// @Tags H5-Chat
// @Param sender_type query string false "sender type (agent/system)"
// @Success 200 {object} model.Response
// @Router /api/v1/chat/unread-count [get]
func (a *ChatApi) GetUnreadCount(c *gin.Context) {
	sessionID := c.GetInt64("session_id")
	senderType := c.DefaultQuery("sender_type", "agent")

	count, err := ServiceGroupApp.ChatService.GetUnreadCount(c.Request.Context(), sessionID, senderType)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	model.OkWithData(gin.H{"count": count}, c)
}

// CloseSession godoc
// @Summary Close current session (called by user)
// @Tags H5-Chat
// @Success 200 {object} model.Response
// @Router /api/v1/chat/close [post]
func (a *ChatApi) CloseSession(c *gin.Context) {
	sessionID := c.GetInt64("session_id")

	err := ServiceGroupApp.ChatService.CloseSession(c.Request.Context(), sessionID)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	model.OkWithMessage("success", c)
}
