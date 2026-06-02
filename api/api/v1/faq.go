package v1

import (
	"github.com/cngamesdk/live-chat/api/model"
	"github.com/gin-gonic/gin"
)

type FaqApi struct{}

type QueryFaqRequest struct {
	Keyword  string `form:"keyword" json:"keyword"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

// QueryFaq godoc
// @Summary Query FAQ by keyword
// @Tags H5-FAQ
// @Param keyword query string true "search keyword"
// @Param page query int false "page"
// @Param pageSize query int false "pageSize"
// @Success 200 {object} model.Response
// @Router /api/v1/chat/faq [get]
func (a *FaqApi) QueryFaq(c *gin.Context) {
	productID := c.GetInt64("product_id")
	var req QueryFaqRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		model.FailWithMessage("invalid params", c)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := ServiceGroupApp.FaqService.Search(c.Request.Context(), productID, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		model.FailWithMessage(err.Error(), c)
		return
	}

	model.OkWithPage(list, total, req.Page, req.PageSize, c)
}
