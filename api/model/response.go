package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	CodeSuccess      = 0
	CodeError        = 1
	CodeParamError   = 400
	CodeUnauthorized = 401
	CodeNotFound     = 404
	CodeServerError  = 500
)

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Ok(c *gin.Context) {
	Result(CodeSuccess, "success", nil, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(CodeSuccess, "success", data, c)
}

func OkWithMessage(msg string, c *gin.Context) {
	Result(CodeSuccess, msg, nil, c)
}

func Fail(c *gin.Context) {
	Result(CodeError, "error", nil, c)
}

func FailWithMessage(msg string, c *gin.Context) {
	Result(CodeError, msg, nil, c)
}

func FailWithCode(code int, msg string, c *gin.Context) {
	Result(code, msg, nil, c)
}

func FailWithData(code int, msg string, data interface{}, c *gin.Context) {
	Result(code, msg, data, c)
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func OkWithPage(list interface{}, total int64, page, pageSize int, c *gin.Context) {
	OkWithData(PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, c)
}
