package controllers

import (
	// "bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	// "bluebell/pkg/jwt"

	// "errors"
	"fmt"
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignupHandler 注册路由
// @Summary 用户注册请求接口
// @Description user signup，输入账号、密码、确认密码，完成注册
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamSignup true "注册参数(账号、密码、确认密码)"
// @Security ApiKeyAuth
// @Success 200 {object} Response "注册成功"
// @Router /signup [post]
func SignupHandler(ctx *gin.Context) {
	// 1.获取参数
	var p models.ParamSignup
	if err := ctx.ShouldBindJSON(&p); err != nil {
		// 类型断言进行转化
		err2, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("绑定参数错误", zap.Error(err))
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(err2.Translate(trans)))
		//    ctx.JSON(http.StatusOK,&gin.H{
		// 	"msg":removeTopStruct(err2.Translate(trans)),
		//    })
		return
	}
	// if p.Username == "" || p.Password == "" || p.AcPassword == "" || p.Password != p.AcPassword {
	// 	zap.L().Error("绑定参数错误")
	// 	ctx.JSON(http.StatusOK, &gin.H{
	// 		"msg": "绑定参数错误",
	// 	})
	// 	return
	// }
	// 2.业务处理
	if err := logic.Signup(&p); err != nil {
		zap.L().Error("注册的错误信息", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeServerBusy, fmt.Sprintf("%v", err))
		// ctx.JSON(http.StatusOK, &gin.H{
		// 	"msg": fmt.Sprintf("%v",err),
		// })
		return
	}
	// 3.返回结果
	ResponseSuccess(ctx, "注册")
	// ctx.JSON(http.StatusOK, &gin.H{
	// 	"msg": "success",
	// })
}

// LoginHandler 登录
// @Summary 用户登录请求接口
// @Description user login，输入账号和密码，得到token继而开放访问其他页面
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLogin false "登录参数(账号和密码)"
// @Security ApiKeyAuth
// @Success 200 {object} Response "登录成功"
// @Router /login [post]
func LoginHandler(ctx *gin.Context) {
	// 1.获取参数
	var p models.ParamLogin
	if err := ctx.ShouldBindJSON(&p); err != nil {
		err2, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			zap.L().Error("绑定参数错误！", zap.Error(err))
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(err2.Translate(trans))) //翻译Validation类型的错误
		return

	}
	// 2.业务处理
	if err := logic.Login(&p); err != nil {
		zap.L().Error("业务处理错误！", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeServerBusy, fmt.Sprintf("%v", err))
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, fmt.Sprintf("Token:%v", logic.UserToken))
}
