package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"Go_skeleton/pkg/config"
	"Go_skeleton/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// RestServer 是对 Gin 的一层包装
type RestServer struct {
	Engine *gin.Engine  // 核心 Gin 引擎
	port   int          // 监听端口
	server *http.Server // 标准库的 http.Server
}

// NewRestServer 像个工厂，根据配置生产一个 Server
func NewRestServer(cfg config.ServerConfig) *RestServer {
	// 1. 设置 Gin 的运行模式 (Debug/Release)
	gin.SetMode(cfg.Mode)

	// 2. 创建 Gin 实例
	r := gin.New()

	// 注册中间件

	// 链路追踪 (Trace)
	r.Use(otelgin.Middleware(cfg.Name))

	//崩溃恢复 (Recovery)
	r.Use(gin.Recovery())

	// 日志记录 (Logger)
	r.Use(ginLogger())

	//  4. 注册基础路由

	// D. 监控指标 (Metrics)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return &RestServer{
		Engine: r,
		port:   cfg.Port,
	}
}

// Start 启动服务
func (s *RestServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.Engine,
	}
	logger.Log.Info("HTTP 服务正在启动...", zap.Int("port", s.port))

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP 服务启动失败: %w", err)
	}
	return nil
}

// Stop 优雅停止
func (s *RestServer) Stop(ctx context.Context) error {
	logger.Log.Info("HTTP 服务正在停止...")
	return s.server.Shutdown(ctx)
}

// ginLogger 中间件
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		logger.Log.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("cost", cost),
		)
	}
}
