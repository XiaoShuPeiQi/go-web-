package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	res := &Response{
		Code: code,
		Msg:  code.getMsg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, res)
}
func ResponseSuccess(c *gin.Context, data interface{}) {
	res := &Response{
		Code: CodeSuccess,
		Msg:  CodeSuccess.getMsg(),
		Data: data,
	}
	c.JSON(http.StatusOK, res)
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	res := &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, res)
}
