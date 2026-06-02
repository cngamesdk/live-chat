package service

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/model"
	"gorm.io/gorm"
)

type AgentService struct{}

var roundRobinCounter int64

func (s *AgentService) GetOnlineAgents(ctx context.Context, productID int64) ([]model.Agent, error) {
	var agents []model.Agent
	err := global.GVA_DB.Where("product_id = ? AND status = ?", productID, "online").Find(&agents).Error
	return agents, err
}

// AssignAgent assigns a session to the best available agent based on strategy.
func (s *AgentService) AssignAgent(ctx context.Context, productID int64) (*model.Agent, error) {
	var agents []model.Agent
	err := global.GVA_DB.Where("product_id = ? AND status = ? AND current_sessions < max_concurrent",
		productID, "online").Find(&agents).Error
	if err != nil {
		return nil, err
	}

	if len(agents) == 0 {
		return nil, nil // No available agents
	}

	strategy := global.GVA_CONFIG.Chat.AgentAssignment

	var selected *model.Agent
	switch strategy {
	case "least_loaded":
		// Find agent with fewest current sessions
		minSessions := int(^uint(0) >> 1)
		for i := range agents {
			if agents[i].CurrentSessions < minSessions {
				minSessions = agents[i].CurrentSessions
				selected = &agents[i]
			}
		}
	default: // round_robin
		idx := atomic.AddInt64(&roundRobinCounter, 1) % int64(len(agents))
		selected = &agents[idx]
	}

	if selected == nil && len(agents) > 0 {
		selected = &agents[0]
	}

	if selected != nil {
		// Increment agent's current sessions
		global.GVA_DB.Model(selected).UpdateColumns(map[string]interface{}{
			"current_sessions": selected.CurrentSessions + 1,
			"total_served":     selected.TotalServed + 1,
		})
	}

	return selected, nil
}

func (s *AgentService) ReleaseAgent(ctx context.Context, agentID int64) error {
	return global.GVA_DB.Model(&model.Agent{}).Where("id = ? AND current_sessions > 0", agentID).
		UpdateColumn("current_sessions", gorm.Expr("current_sessions - 1")).Error
}

func (s *AgentService) GoOnline(ctx context.Context, userID, productID int64) error {
	var agent model.Agent
	err := global.GVA_DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&agent).Error
	now := time.Now()
	if err != nil {
		agent = model.Agent{
			UserID:       userID,
			ProductID:    productID,
			Status:       "online",
			MaxConcurrent: 5,
			LastOnlineAt: &now,
		}
		return global.GVA_DB.Create(&agent).Error
	}
	return global.GVA_DB.Model(&agent).Updates(map[string]interface{}{
		"status":         "online",
		"last_online_at": now,
	}).Error
}

func (s *AgentService) GoOffline(ctx context.Context, userID, productID int64) error {
	return global.GVA_DB.Model(&model.Agent{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("status", "offline").Error
}

func (s *AgentService) UpdateConfig(ctx context.Context, id int64, maxConcurrent int) error {
	return global.GVA_DB.Model(&model.Agent{}).Where("id = ?", id).
		Update("max_concurrent", maxConcurrent).Error
}

func (s *AgentService) List(ctx context.Context, productID int64, page, pageSize int) ([]model.Agent, int64, error) {
	var list []model.Agent
	var total int64
	db := global.GVA_DB.Model(&model.Agent{})
	if productID > 0 {
		db = db.Where("product_id = ?", productID)
	}
	db.Count(&total)
	err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}
