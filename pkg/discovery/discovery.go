package discovery

import (
	"context"
	"fmt"
	"time"

	"Go_skeleton/pkg/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

// EtcdResolver 负责从 Etcd 监听服务变化，并告诉 gRPC
type EtcdResolver struct {
	client *clientv3.Client
}

// NewEtcdResolver 创建解析器 Builder
func NewEtcdResolver(addr []string) (*EtcdResolver, error) {
	cfg := clientv3.Config{Endpoints: addr, DialTimeout: 5 * time.Second}
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return &EtcdResolver{client: cli}, nil
}

// Build 是 resolver.Builder 接口的方法
// gRPC 在拨号时调用它，它负责启动解析逻辑
func (r *EtcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// target.Endpoint 就是服务名，比如 "user-service"
	serviceName := target.Endpoint()
	// 前缀，跟注册时保持一致
	prefix := fmt.Sprintf("/micro/services/%s/", serviceName)

	// 1. 刚启动时，先主动拉取一次现有的 IP
	r.updateState(context.Background(), prefix, cc)

	// 2. 启动监听 (Watch)，如果有变化自动通知 gRPC
	go r.watch(prefix, cc)

	return r, nil
}

// Scheme 定义我们的协议头，即 grpc.Dial("etcd:///...")
func (r *EtcdResolver) Scheme() string {
	return "etcd"
}

// --- 下面是补全的方法，解决报错的核心 ---

// ResolveNow 是 resolver.Resolver 接口的方法
// 当 gRPC 觉得连接有问题时，会调用它让我们马上再刷新一下
func (r *EtcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	// 这里可以留空，或者打个日志，因为我们有 Watch 机制自动处理
	logger.Log.Debug("ResolveNow 被调用，但我们通过 Watch 机制更新，所以忽略")
}

// Close 是 resolver.Resolver 接口的方法
// 当 gRPC 断开连接时调用
func (r *EtcdResolver) Close() {
	logger.Log.Info("EtcdResolver 关闭")
	// 这里可以根据需要关闭 etcd client，但通常 Resolver 是全局复用的，不一定要关
}

// --- 内部辅助逻辑 ---

// watch 监听变化
func (r *EtcdResolver) watch(prefix string, cc resolver.ClientConn) {
	// 监听前缀的变化
	rch := r.client.Watch(context.Background(), prefix, clientv3.WithPrefix())

	for n := range rch {
		for _, _ = range n.Events { // 【修复】把 ev 改成 _，解决 Unused variable 报错
			// 只要有事件（新增或删除），就重新拉取一次最新的列表
			logger.Log.Info("检测到服务地址变化，正在更新...", zap.String("prefix", prefix))
			r.updateState(context.Background(), prefix, cc)
		}
	}
}

// updateState 封装了“去 Etcd 查 IP 并告诉 gRPC”的逻辑
// 这样 Build 和 watch 都可以调用它，避免代码重复
func (r *EtcdResolver) updateState(ctx context.Context, prefix string, cc resolver.ClientConn) {
	// 1. 去 Etcd 查
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Log.Error("Etcd 获取服务地址失败", zap.Error(err))
		return
	}

	// 2. 解析出 IP 列表
	var addrs []resolver.Address
	for _, kv := range resp.Kvs {
		// Value 里存的是 "127.0.0.1:8080"
		addr := string(kv.Value)
		addrs = append(addrs, resolver.Address{Addr: addr})
		logger.Log.Debug("发现服务节点", zap.String("addr", addr))
	}

	// 3. 告诉 gRPC：现在这些 IP 是可用的
	// gRPC 内部会拿着这些 IP 做负载均衡
	err = cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		logger.Log.Error("更新 gRPC 状态失败", zap.Error(err))
	}
}
