package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建帖子接口：输入标题、内容、所属社区id
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object body models.PostSwagger false "创建帖子参数"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /post [post]
func CreatePostHandler(ctx *gin.Context) {
	// 1.获取帖子参数
	p := new(models.Post)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("ShouldBindJSON绑定结构体出错", zap.Error(err))
		return
	}
	// 2.写入数据库
	if err := logic.CreatePost(ctx, p); err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, "创建post成功")
}

// GetPostHandler 根据id获取帖子信息
// @Summary 根据id获取帖子
// @Description 输入帖子id，获取帖子信息
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /post/{id} [get]
func GetPostHandler(ctx *gin.Context) {
	//	1.获取帖子id
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	//	2.查数据库得到数据
	data, err := logic.GetPost(id)
	if err != nil {
		zap.L().Error("logic.GetPost(id) wrong", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//	3.返回响应
	ResponseSuccess(ctx, data)

}

// GetPostListHandler 分页获取帖子列表
// @Summary 分页获取帖子列表
// @Description 输入page和size，获取帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param size query string false "每页大小"
// @Param page query string false "第几页？"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /posts [get]
func GetPostListHandler(ctx *gin.Context) {
	//  获取参数page和size
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 100
	}
	//	获取帖子列表
	PostList, err := logic.GetPostList(size, page)
	if err != nil {
		zap.L().Error("GetPostList wrong", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeServerBusy, err)
		return
	}
	//	返回响应
	ResponseSuccess(ctx, PostList)

}

// GetPostListHandler2 新版社区列表，支持按发帖时间/帖子分数排序
// @Summary 新版社区列表，支持按发帖时间/帖子分数排序
// @Description 输入page、size、order，获取帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param size query string false "每页大小"
// @Param page query string false "第几页？"
// @Param order query string false "排序方式(time/score)"
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Router /posts2 [get]
func GetPostListHandler2(ctx *gin.Context) {
	//1.获取参数  /pages2?size=10&page=1&order=time
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: "score",
	}
	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 ShouldBindQuery wrong", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
	}
	fmt.Println(p)
	//2.交给logic，返回列表
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList wrong", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
	}
	//3.返回响应
	ResponseSuccess(ctx, data)
}

// PinHandler 处理ping请求
// @Summary Token有效性检测
// @Description 这个接口用于ping服务器，返回"pong"消息。
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Router /ping [post]
func PinHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "pong",
	})
}
