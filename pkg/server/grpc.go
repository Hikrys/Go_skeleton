package server

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"Go_skeleton/pkg/config"
	"Go_skeleton/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RpcServer 包装 gRPC Server
type RpcServer struct {
	Server *grpc.Server // 真正的 gRPC Server 对象
	port   int
}

func NewRpcServer(cfg config.ServerConfig) *RpcServer {
	// 这里是核心：添加拦截器 (Middleware)
	// 如果你有多个拦截器，可以使用 grpc.ChainUnaryInterceptor 串联起来
	opts := []grpc.ServerOption{
		// 1. OTel 注入 (放在 StatsHandler 里)
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		//5 2.拦截器链
		grpc.ChainUnaryInterceptor(
			// Recovery 拦截器：防止 panic 崩溃
			recoveryInterceptor,
			// Timeout 拦截器：防止请求处理时间过长
			timeoutInterceptor(5*time.Second), // 假设默认超时 5 秒
			// 如果还有 Auth, Trace 都在这里加
		),
	}

	s := grpc.NewServer(opts...)

	return &RpcServer{
		Server: s,
		port:   cfg.Port + 1000, // 简单起见，RPC 端口比 HTTP 大 1000，比如 8080 -> 9080
	}
}

// Start 启动 RPC 服务
func (s *RpcServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("RPC 端口监听失败: %w", err)
	}

	logger.Log.Info("gRPC 服务正在启动...", zap.Int("port", s.port))

	// Serve 也是阻塞的
	return s.Server.Serve(lis)
}

// Stop 优雅停止
func (s *RpcServer) Stop() {
	logger.Log.Info("gRPC 服务正在停止...")
	// GracefulStop 会等待所有正在进行的 RPC 调用结束
	s.Server.GracefulStop()
}

// --- 拦截器实现细节 (Interceptor) ---

// recoveryInterceptor 就像一道保险丝
func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// defer 作为兜底，panic 发生时会触发
	defer func() {
		if e := recover(); e != nil {
			// 打印堆栈信息，方便查错
			logger.Log.Error("RPC 服务发生 Panic!",
				zap.Any("error", e),
				zap.String("stack", string(debug.Stack())),
			)
			// 返回给客户端一个明确的错误
			err = status.Errorf(codes.Internal, "服务内部错误")
		}
	}()
	// 继续执行真正的业务逻辑
	return handler(ctx, req)
}

// timeoutInterceptor 给每个请求加上倒计时
func timeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 派生出一个带超时的 context
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel() // 记得释放

		// 传进去带超时的 ctx
		// 如果 handler 处理太慢，ctx.Done() 就会先结束
		return handler(ctx, req)
	}
}
