package ws

import (
	"time"

	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/cngamesdk/live-chat/api/service"
	"go.uber.org/zap"
)

func HandleClientMessage(c *Client, msg *Message) {
	switch msg.Type {
	case "message":
		handleChatMessage(c, msg)
	case "ping":
		c.Send <- []byte(`{"type":"pong","timestamp":0}`)
	}
}

func handleChatMessage(c *Client, msg *Message) {
	ctx := c.ctx()
	msgType := msg.MsgType
	if msgType == "" {
		msgType = "text"
	}

	chatSvc := new(service.ChatService)
	savedMsg, faqMatch, err := chatSvc.ProcessUserMessage(
		ctx,
		c.SessionID,
		msg.Content,
		msgType,
		msg.AttachmentURL,
	)
	if err != nil {
		logger.Error(ctx, "process message error", zap.Error(err))
		return
	}

	GlobalHub.SendToSession(c.SessionID, &Message{
		Type:          "message",
		Content:       savedMsg.Content,
		MsgType:       savedMsg.MsgType,
		SenderType:    savedMsg.SenderType,
		SenderID:      savedMsg.SenderID,
		SenderName:    savedMsg.SenderName,
		AttachmentURL: savedMsg.AttachmentURL,
		MessageID:     savedMsg.ID,
		IsRead:        savedMsg.IsRead,
		Timestamp:     time.Now().UnixMilli(),
	})

	if faqMatch != nil {
		GlobalHub.SendToSession(c.SessionID, &Message{
			Type:       "message",
			Content:    faqMatch.Faq.Answer,
			MsgType:    "text",
			SenderType: "system",
			SenderName: "智能客服",
			IsFaqReply: true,
			Timestamp:  time.Now().UnixMilli(),
		})

		session, _ := chatSvc.GetSession(ctx, c.SessionID)
		if session == nil || session.Status == "waiting" {
			tryAssignAgent(c)
		}
		return
	}

	if session, _ := chatSvc.GetSession(ctx, c.SessionID); session != nil {
		if session.Status == "waiting" && session.AgentID == nil {
			tryAssignAgent(c)
		}
	}
}

func tryAssignAgent(c *Client) {
	ctx := c.ctx()
	agentSvc := new(service.AgentService)
	agent, err := agentSvc.AssignAgent(ctx, c.ProductID)
	if err != nil || agent == nil {
		GlobalHub.SendSystemMessage(c.SessionID, "当前没有客服在线，请耐心等待...")
		return
	}

	chatSvc := new(service.ChatService)
	chatSvc.AssignAgent(ctx, c.SessionID, agent.ID)

	GlobalHub.AgentConnected(c.SessionID, "客服"+string(rune(agent.ID+'0')))
	GlobalHub.SendSystemMessage(c.SessionID, "已为您接通人工客服，请描述您的问题。")
}
