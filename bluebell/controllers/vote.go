package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// VoteForPostHandler 给帖子投票
// @Summary 给帖子投票
// @Description 帖子id(表明哪个帖子)、投什么票(1：赞成 0：不投票 -1：反对)
// @Tags 投票相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object body models.ParamVoteForPost true "投票参数()"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /vote [post]
func VoteForPostHandler(ctx *gin.Context) {
	//	1.获取参数
	p := new(models.ParamVoteForPost)
	if err := ctx.ShouldBindJSON(p); err != nil {
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("ShouldBindJSON wrong", zap.Error(err))
			ResponseError(ctx, CodeServerBusy)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(err2.Translate(trans)))
		return
	}
	//	2.给logic层做业务处理
	if err := logic.VoteForPost(ctx, p); err != nil {
		zap.L().Error("logic.VoteForPost wrong", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//	3.返回响应， 成功与否
	ResponseSuccess(ctx, p)
}
