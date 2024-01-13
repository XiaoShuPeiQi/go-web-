package controllers

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommunityHandler get社区表
// @Summary get社区表
// @Description 访问此接口，返回社区信息切片
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /community [get]
func CommunityHandler(ctx *gin.Context) {
	// 响应从社区表中拿到的数据
	list, err := logic.GetCommunityList()
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
	}
	ResponseSuccess(ctx, list)
}

// CommunityDetailHandler 根据id返回社区详情
// @Summary 根据id返回社区详情
// @Description 访问此接口，返回指定id的社区信息
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "社区ID(1-4)"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /community/{id} [get]
func CommunityDetailHandler(ctx *gin.Context) {
	// 拿到id
	idstr := ctx.Param("id")
	// 传递给logic层， 返回数据
	id, _ := strconv.ParseInt(idstr, 10, 64)
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		ResponseErrorWithMsg(ctx, CodeServerBusy, "请求参数错误")
		return
	}
	ResponseSuccess(ctx, data)
}
