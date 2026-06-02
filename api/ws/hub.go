package ws

import (
	"encoding/json"
	"sync"
)

// Hub maintains the set of active clients grouped by session.
type Hub struct {
	mu       sync.RWMutex
	sessions map[int64]map[*Client]bool // sessionID -> clients

	// Message channels for agent notifications
	broadcast chan *Message

	// Register/Unregister
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	Type          string   `json:"type"`
	Content       string   `json:"content,omitempty"`
	MsgType       string   `json:"msg_type,omitempty"`
	SenderType    string   `json:"sender_type,omitempty"`
	SenderID      string   `json:"sender_id,omitempty"`
	SenderName    string   `json:"sender_name,omitempty"`
	AgentName     string   `json:"agent_name,omitempty"`
	AttachmentURL string   `json:"attachment_url,omitempty"`
	IsFaqReply    bool     `json:"is_faq_reply,omitempty"`
	IsRead        bool     `json:"is_read,omitempty"`
	MessageID     int64    `json:"message_id,omitempty"`
	QueuePosition int      `json:"queue_position,omitempty"`
	SessionID     int64    `json:"session_id,omitempty"`
	Timestamp     int64    `json:"timestamp"`
	Data          any      `json:"data,omitempty"`
}

var GlobalHub *Hub

func NewHub() *Hub {
	return &Hub{
		sessions:   make(map[int64]map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.sessions[client.SessionID]; !ok {
				h.sessions[client.SessionID] = make(map[*Client]bool)
			}
			h.sessions[client.SessionID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.sessions[client.SessionID]; ok {
				delete(clients, client)
				close(client.Send)
				if len(clients) == 0 {
					delete(h.sessions, client.SessionID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.sessions[msg.SessionID]; ok {
				data, _ := json.Marshal(msg)
				for client := range clients {
					select {
					case client.Send <- data:
					default:
						go func(c *Client) { h.unregister <- c }(client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToSession(sessionID int64, msg *Message) {
	msg.SessionID = sessionID
	h.broadcast <- msg
}

func (h *Hub) SendSystemMessage(sessionID int64, content string) {
	h.SendToSession(sessionID, &Message{
		Type:       "message",
		Content:    content,
		MsgType:    "system",
		SenderType: "system",
		SenderName: "系统",
	})
}

func (h *Hub) QueuePosition(sessionID int64, position int) {
	h.SendToSession(sessionID, &Message{
		Type:          "queue",
		QueuePosition: position,
	})
}

func (h *Hub) AgentConnected(sessionID int64, agentName string) {
	h.SendToSession(sessionID, &Message{
		Type:       "connected",
		AgentName:  agentName,
		SenderName: "系统",
	})
}

func (h *Hub) SessionClosed(sessionID int64) {
	h.SendToSession(sessionID, &Message{
		Type: "closed",
	})
}

func (h *Hub) GetActiveSessionCount(sessionID int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.sessions[sessionID]; ok {
		return len(clients)
	}
	return 0
}
