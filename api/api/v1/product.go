package v1

import (
	"github.com/cngamesdk/live-chat/api/model"
	"github.com/gin-gonic/gin"
)

type ProductApi struct{}

// GetProductInfo godoc
// @Summary Get product info by code
// @Tags H5-Product
// @Param code path string true "product code"
// @Success 200 {object} model.Response
// @Router /api/v1/product/{code}/info [get]
func (a *ProductApi) GetProductInfo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		model.FailWithMessage("product code required", c)
		return
	}

	product, err := ServiceGroupApp.ProductService.GetByCode(c.Request.Context(), code)
	if err != nil {
		model.FailWithMessage("product not found", c)
		return
	}

	model.OkWithData(gin.H{
		"product_code":    product.ProductCode,
		"name":            product.Name,
		"logo":            product.Logo,
		"welcome_title":   product.WelcomeTitle,
		"welcome_message": product.WelcomeMessage,
	}, c)
}
