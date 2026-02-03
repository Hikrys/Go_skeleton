package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Go_skeleton/pkg/logger"
	"Go_skeleton/pkg/server"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// App 是整个程序的管家
type App struct {
	RestServer *server.RestServer
	RpcServer  *server.RpcServer
}

// NewApp 把所有的 Server 组装起来
func NewApp(rest *server.RestServer, rpc *server.RpcServer) *App {
	return &App{
		RestServer: rest,
		RpcServer:  rpc,
	}
}

// Run 启动所有
func (a *App) Run() error {
	// 创建一个带取消功能的 Context
	// errgroup 需要这个 context 来通知大家“该停了”
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建 errgroup
	g, ctx := errgroup.WithContext(ctx)

	// ---  1: 启动 HTTP 服务 ---
	g.Go(func() error {
		// 启动服务
		err := a.RestServer.Start()
		if err != nil {
			logger.Log.Error("Rest Server 挂了", zap.Error(err))
		}
		// 如果 Start 返回（说明服务挂了），我们要调用 cancel 取消其他任务
		cancel()
		return err
	})

	// ---  2: 启动 gRPC 服务 ---
	g.Go(func() error {
		err := a.RpcServer.Start()
		if err != nil {
			logger.Log.Error("Rpc Server 挂了", zap.Error(err))
		}
		cancel()
		return err
	})

	// ---3: 监听系统的退出信号 (Ctrl+C) ---
	g.Go(func() error {
		// 创建一个通道，专门接收信号
		quit := make(chan os.Signal, 1)
		// 监听 SIGINT (Ctrl+C) 和 SIGTERM (Docker/K8s 停止命令)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			// 如果 context 已经取消了（说明别的服务挂了），那我们直接退出
			return ctx.Err()
		case sig := <-quit:
			// 如果收到了系统信号
			logger.Log.Info("接收到退出信号，准备优雅关闭...", zap.String("signal", sig.String()))
			// 这里调用 cancel() 告诉 errgroup ,让errgroup里面所有的任务停止
			cancel()
		}
		return nil
	})

	// ---  4: 负责执行具体停止逻辑 ---
	// 上面的任务只是“通知”要停，这个任务负责“执行”停止
	g.Go(func() error {
		// 等待 context 被取消（也就是收到了停止通知）
		<-ctx.Done()

		logger.Log.Info("正在关闭所有服务...")

		// 1. 关闭 HTTP (给它 5 秒钟时间处理完剩下的请求)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := a.RestServer.Stop(shutdownCtx); err != nil {
			logger.Log.Error("HTTP 关闭出错", zap.Error(err))
		}

		// 2. 关闭 RPC
		a.RpcServer.Stop()

		logger.Log.Info("所有服务已关闭，再见！")
		return nil
	})

	// Wait 等待所有任务结束
	// 只要有一个任务返回 error，Wait 就会返回那个 error
	return g.Wait()
}
