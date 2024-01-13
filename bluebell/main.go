package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"fmt"

	"go.uber.org/zap"
)

// @title bluebell 项目接口文档
// @version 2.0
// @description Go web 开发进阶项目实战 --> 源自李文周
// @contact.name XiaoShuPeiQi
// @contact.email 1446596766@qq.com
// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080/v1
func main() {
	// 1. 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Println("加载配置文件错误")
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Println("初始化日志错误")
		return
	}
	// 3. 初始化mysql
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		zap.L().Error("初始化mysql错误", zap.Error(err))
		return
	}
	defer mysql.Close()
	// 4. 初始化redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Error("初始化redis错误", zap.Error(err))
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, int64(settings.Conf.MachineID)); err != nil {
		zap.L().Error("初始化雪花算法错误", zap.Error(err))
		return

	}

	// 5. 注册路由
	r := routes.Setup()
	// 6. 启动服务
	r.Run(":8080")
}
