package routes

import (
	"bluebell/controllers"
	_ "bluebell/docs" // 千万不要忘了导入把你上一步生成的docs
	"bluebell/logger"
	"bluebell/middlewares"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"time"
)

func Setup() *gin.Engine {

	// 初始化翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Println("初始化翻译器错误")
	}
	r := gin.New()

	//使用中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.Use(middlewares.CORSMiddleware())
	//r.LoadHTMLFiles("./templates/index.html")
	//r.Static("/static", "./static")
	////注册路由业务
	//r.GET("/", func(ctx *gin.Context) {
	//	log.Println("Received a request to /")
	//	ctx.HTML(http.StatusOK, "index.html", nil)
	//})
	v1 := r.Group("/v1")
	{
		// 注册
		v1.POST("/signup", controllers.SignupHandler)
		// 登录
		v1.POST("/login", controllers.LoginHandler)
		// 登录过后想访问的界面
		v1.Use(middlewares.JWTAuthMiddleware())
		{

			v1.GET("/community", controllers.CommunityHandler)
			v1.GET("/community/:id", controllers.CommunityDetailHandler)

			v1.POST("/post", controllers.CreatePostHandler)
			v1.GET("/post/:id", controllers.GetPostHandler)
			v1.GET("/posts", controllers.GetPostListHandler)
			v1.GET("/posts2", controllers.GetPostListHandler2)

			v1.POST("/vote", controllers.VoteForPostHandler)

			v1.Use(middlewares.RateLimitMiddleware(2*time.Second, 2))
			v1.POST("/ping", controllers.PinHandler)

		}
		r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}
	pprof.Register(r) //注册性能分析pprof
	return r
}
