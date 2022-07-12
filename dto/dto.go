package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseFormat 响应格式
type ResponseFormat struct {
	Code uint64      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// FailResponse 失败响应
func FailResponse(context *gin.Context, status int, msg string) {
	format := ResponseFormat{
		Code: 1,
		Msg:  msg,
	}

	context.JSON(status, format)
}

// SuccessResponse 成功响应
func SuccessResponse(context *gin.Context, data interface{}) {
	resp := ResponseFormat{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
	context.JSON(http.StatusOK, resp)
}
