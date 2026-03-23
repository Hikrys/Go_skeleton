package registry

import (
	"context"
	"fmt"
	"time"

	"Go_skeleton/pkg/logger"
	clientv3 "go.etcd.io/etcd/client/v3" // 这是官方的 Etcd 客户端库
	"go.uber.org/zap"
)

// Registry 定义一个通用的接口
type Registry interface {
	Register(ctx context.Context, serviceName string, host string, port int) error
	Unregister(ctx context.Context, serviceName string) error
}

// EtcdRegistry 是具体的实现
type EtcdRegistry struct {
	client  *clientv3.Client // Etcd 的连接客户端
	leaseID clientv3.LeaseID // 租约 ID
}

// NewEtcdRegistry 创建一个注册器
// addr 是 Etcd 的地址，比如 ["127.0.0.1:2379"]
func NewEtcdRegistry(addr []string) (*EtcdRegistry, error) {
	// 配置 Etcd 客户端
	cfg := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("无法连接 Etcd: %w", err)
	}

	return &EtcdRegistry{
		client: client,
	}, nil
}

func (r *EtcdRegistry) Register(ctx context.Context, serviceName string, host string, port int) error {
	// 创建一个租约 (Lease)
	// 5 表示 5 秒。如果 5 秒没续租，Etcd 认为服务消失
	leaseResp, err := r.client.Grant(ctx, 5)
	if err != nil {
		return fmt.Errorf("创建租约失败: %w", err)
	}
	r.leaseID = leaseResp.ID

	// 生成存到 Etcd 里的 Key 和 Value
	// Key: /micro/services/user-service/127.0.0.1:8080
	// 这样设计是为了方便前缀搜索
	key := fmt.Sprintf("/micro/services/%s/%s:%d", serviceName, host, port)
	value := fmt.Sprintf("%s:%d", host, port)

	// 把数据写入 Etcd，并绑定租约
	// WithLease 意思是：这个 Key 的寿命跟这个租约绑定
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("注册服务信息失败: %w", err)
	}

	// 开启自动续租
	// 这是一个长连接，会在后台一直发心跳
	keepAliveCh, err := r.client.KeepAlive(ctx, r.leaseID)
	if err != nil {
		return fmt.Errorf("启动心跳失败: %w", err)
	}

	// 启动一个协程去处理心跳的响应
	// 虽然我们不需要对响应做什么，但必须把 channel 里的数据读出来，不然 channel 满了会阻塞
	go func() {
		for {
			select {
			case _, ok := <-keepAliveCh:
				if !ok {
					logger.Log.Warn("Etcd 心跳通道关闭了，可能是 Etcd 挂了或者网络断了")
					return
				}
			case <-ctx.Done():
				logger.Log.Info("停止心跳续租")
				return
			}
		}
	}()

	logger.Log.Info("服务注册成功", zap.String("key", key), zap.Int64("leaseID", int64(r.leaseID)))
	return nil
}

// Unregister 注销
func (r *EtcdRegistry) Unregister(ctx context.Context, serviceName string) error {
	// 直接撤销租约
	// 租约一没，绑定的 Key 自动删除，不需要手动 Delete Key
	if _, err := r.client.Revoke(ctx, r.leaseID); err != nil {
		return fmt.Errorf("注销服务失败: %w", err)
	}
	logger.Log.Info("服务注销成功")
	return nil
}
