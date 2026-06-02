package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/model"
	"gorm.io/gorm"
)

type ProductService struct{}

func (s *ProductService) GetByCode(ctx context.Context, code string) (*model.Product, error) {
	cacheKey := global.GVA_CONFIG.Redis.Prefix + "product:" + code

	// Try Redis cache
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var p model.Product
			if json.Unmarshal([]byte(cached), &p) == nil {
				return &p, nil
			}
		}
	}

	var p model.Product
	err := global.GVA_DB.Where("product_code = ? AND status = ?", code, "normal").First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Cache for 5 minutes
	if global.GVA_REDIS != nil {
		data, _ := json.Marshal(p)
		global.GVA_REDIS.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return &p, nil
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	return global.GVA_DB.Create(p).Error
}

func (s *ProductService) Update(ctx context.Context, p *model.Product) error {
	return global.GVA_DB.Save(p).Error
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return global.GVA_DB.Delete(&model.Product{}, id).Error
}

func (s *ProductService) List(ctx context.Context, page, pageSize int) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64
	db := global.GVA_DB.Model(&model.Product{})
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&list).Error
	return list, total, err
}
