package service

import (
	"context"
	"sort"
	"strings"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/model"
)

type FaqService struct{}

type FaqMatch struct {
	Faq   model.Faq `json:"faq"`
	Score float64   `json:"score"`
}

func (s *FaqService) Search(ctx context.Context, productID int64, keyword string, page, pageSize int) ([]model.Faq, int64, error) {
	var list []model.Faq
	var total int64
	db := global.GVA_DB.Model(&model.Faq{}).Where("product_id = ? AND status = ?", productID, "normal")

	if keyword != "" {
		db = db.Where("question LIKE ? OR keywords LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	db.Count(&total)
	err := db.Order("priority DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// MatchBest finds the best matching FAQ for a user message.
// Returns the best match if score exceeds threshold, or nil if no match.
func (s *FaqService) MatchBest(ctx context.Context, productID int64, message string) *FaqMatch {
	threshold := global.GVA_CONFIG.Chat.FaqMatchThreshold

	var allFaqs []model.Faq
	global.GVA_DB.Where("product_id = ? AND status = ?", productID, "normal").
		Order("priority DESC").Find(&allFaqs)

	if len(allFaqs) == 0 {
		return nil
	}

	msgLower := strings.ToLower(strings.TrimSpace(message))
	msgWords := tokenize(msgLower)

	var matches []FaqMatch
	for _, faq := range allFaqs {
		score := calculateMatchScore(msgLower, msgWords, faq)
		if score > 0 {
			matches = append(matches, FaqMatch{Faq: faq, Score: score})
		}
	}

	if len(matches) == 0 {
		return nil
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	best := matches[0]
	if best.Score >= threshold {
		// Increment match count asynchronously
		go func() {
			global.GVA_DB.Model(&best.Faq).UpdateColumn("match_count", best.Faq.MatchCount+1)
		}()
		return &best
	}

	return nil
}

func tokenize(text string) []string {
	// Simple tokenization: split by whitespace and common Chinese punctuation
	words := strings.Fields(text)
	var result []string
	for _, w := range words {
		w = strings.Trim(w, "，。！？,.!?")
		if len(w) > 0 {
			result = append(result, w)
		}
	}
	return result
}

func calculateMatchScore(msgLower string, msgWords []string, faq model.Faq) float64 {
	questionLower := strings.ToLower(faq.Question)
	keywordsLower := strings.ToLower(faq.Keywords)

	// Exact question match
	if msgLower == questionLower {
		return 1.0
	}

	// Full question contained in message
	if strings.Contains(msgLower, questionLower) {
		return 0.95
	}

	// Message contained in question
	if strings.Contains(questionLower, msgLower) {
		return 0.9
	}

	// Word-level matching
	questionWords := tokenize(questionLower)
	keywordWords := tokenize(keywordsLower)

	allTargetWords := append(questionWords, keywordWords...)
	if len(allTargetWords) == 0 {
		return 0
	}

	matched := 0
	for _, mw := range msgWords {
		for _, tw := range allTargetWords {
			if mw == tw || strings.Contains(tw, mw) || strings.Contains(mw, tw) {
				matched++
				break
			}
		}
	}

	score := float64(matched) / float64(len(allTargetWords))
	if score > 1.0 {
		score = 1.0
	}
	return score
}
