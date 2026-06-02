package service

import (
	"context"
	"errors"
	"time"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/middleware"
	"github.com/cngamesdk/live-chat/api/model"
	"gorm.io/gorm"
)

type ChatService struct{}

func (s *ChatService) InitSession(ctx context.Context, productCode, userID, userName, source, userIP, userAgent string) (*model.ChatSession, string, error) {
	// Get product
	product, err := new(ProductService).GetByCode(ctx, productCode)
	if err != nil {
		return nil, "", errors.New("product not found or disabled")
	}

	// Check if user has an active session
	var existing model.ChatSession
	err = global.GVA_DB.Where("product_id = ? AND user_id = ? AND status != ?",
		product.ID, userID, "closed").Order("id DESC").First(&existing).Error
	if err == nil {
		// Reuse existing session
		token, tokErr := middleware.GenerateSessionToken(existing.ID, existing.ProductID, existing.UserID)
		if tokErr != nil {
			return nil, "", tokErr
		}
		return &existing, token, nil
	}

	session := model.ChatSession{
		ProductID: product.ID,
		UserID:    userID,
		UserName:  userName,
		Status:    "waiting",
		Source:    source,
		UserIP:    userIP,
		UserAgent: userAgent,
	}

	if err := global.GVA_DB.Create(&session).Error; err != nil {
		return nil, "", err
	}

	// Send welcome message
	welcomeMsg := product.WelcomeMessage
	if welcomeMsg == "" {
		welcomeMsg = "您好，欢迎咨询！请问有什么可以帮助您的？"
	}
	sysMsg := model.ChatMessage{
		SessionID:  session.ID,
		SenderType: "system",
		SenderID:   "system",
		SenderName: "系统",
		Content:    welcomeMsg,
		MsgType:    "system",
	}
	global.GVA_DB.Create(&sysMsg)

	token, err := middleware.GenerateSessionToken(session.ID, session.ProductID, session.UserID)
	if err != nil {
		return nil, "", err
	}

	return &session, token, nil
}

func (s *ChatService) GetHistory(ctx context.Context, sessionID int64, page, pageSize int) ([]model.ChatMessage, int64, error) {
	var list []model.ChatMessage
	var total int64
	db := global.GVA_DB.Model(&model.ChatMessage{}).Where("session_id = ?", sessionID)
	db.Count(&total)
	err := db.Order("created_at ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (s *ChatService) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	return global.GVA_DB.Create(msg).Error
}

func (s *ChatService) GetSession(ctx context.Context, sessionID int64) (*model.ChatSession, error) {
	var session model.ChatSession
	err := global.GVA_DB.First(&session, sessionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (s *ChatService) CloseSession(ctx context.Context, sessionID int64) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.ChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"status":    "closed",
		"closed_at": now,
	}).Error
}

func (s *ChatService) AssignAgent(ctx context.Context, sessionID, agentID int64) error {
	return global.GVA_DB.Model(&model.ChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"agent_id": agentID,
		"status":   "active",
	}).Error
}

func (s *ChatService) ProcessUserMessage(ctx context.Context, sessionID int64, content, msgType, attachmentURL string) (*model.ChatMessage, *FaqMatch, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, nil, err
	}

	userID := session.UserID
	userName := session.UserName
	if userName == "" {
		userName = userID
	}

	// Save user message
	msg := &model.ChatMessage{
		SessionID:     sessionID,
		SenderType:    "user",
		SenderID:      userID,
		SenderName:    userName,
		Content:       content,
		MsgType:       msgType,
		AttachmentURL: attachmentURL,
	}
	if err := s.SaveMessage(ctx, msg); err != nil {
		return nil, nil, err
	}

	// For text messages, try FAQ matching
	if msgType == "text" && session.Status == "waiting" {
		faqMatch := new(FaqService).MatchBest(ctx, session.ProductID, content)
		if faqMatch != nil {
			// Auto-reply with FAQ
			faqMsg := &model.ChatMessage{
				SessionID:  sessionID,
				SenderType: "system",
				SenderID:   "system",
				SenderName: "智能客服",
				Content:    faqMatch.Faq.Answer,
				MsgType:    "text",
				IsFaqReply: true,
				FaqID:      &faqMatch.Faq.ID,
			}
			s.SaveMessage(ctx, faqMsg)
			return msg, faqMatch, nil
		}
	}

	return msg, nil, nil
}

// MarkMessagesAsRead marks messages as read
func (s *ChatService) MarkMessagesAsRead(ctx context.Context, sessionID int64, messageIDs []int64) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.ChatMessage{}).
		Where("session_id = ? AND id IN ?", sessionID, messageIDs).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// GetUnreadCount gets unread message count for a session
func (s *ChatService) GetUnreadCount(ctx context.Context, sessionID int64, senderType string) (int64, error) {
	var count int64
	err := global.GVA_DB.Model(&model.ChatMessage{}).
		Where("session_id = ? AND sender_type = ? AND is_read = ?", sessionID, senderType, false).
		Count(&count).Error
	return count, err
}

// SaveAgentMessage saves agent message and broadcasts via WebSocket
func (s *ChatService) SaveAgentMessage(ctx context.Context, sessionID, agentID int64, agentName, content, msgType string) (*model.ChatMessage, error) {
	msg := &model.ChatMessage{
		SessionID:  sessionID,
		SenderType: "agent",
		SenderID:   string(rune(agentID + '0')),
		SenderName: agentName,
		Content:    content,
		MsgType:    msgType,
	}

	if err := s.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}
