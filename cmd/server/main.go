package main

import (
	"Go_skeleton/internal/service"
	"Go_skeleton/pkg/app"
	"Go_skeleton/pkg/config"
	"Go_skeleton/pkg/logger"
	"Go_skeleton/pkg/trace"
	"flag"
)

var configFile = flag.String("f", "configs/config.yaml", "配置文件路径")

// 定义一个“集装箱”结构体

type AppContainer struct {
	App         *app.App
	UserService *service.UserService
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	logger.InitLogger(cfg.Log)

	// 初始化链路追踪
	cleanupTrace, err := trace.InitTracer(cfg.Server.Name, "http://localhost:14268/api/traces")
	if err != nil {
		panic(err)
	}
	defer cleanupTrace(nil)

	// 依赖注入
	// initApp 现在返回一个集装箱 container
	container, cleanup, err := initApp(cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup() // 调用 cleanup
	// 从集装箱里拿出需要的东西
	// 注册路由
	container.UserService.RegisterRoutes(container.App.RestServer.Engine)

	// 启动
	if err := container.App.Run(); err != nil {
		panic(err)
	}
}
