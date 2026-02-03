# Go_skeleton (Go 微服务脚手架)

🚀 一个基于 Go 语言 (1.20+) 的工业级微服务脚手架，集成了 Gin、gRPC、Etcd、OpenTelemetry、Prometheus 等主流技术栈。旨在帮助开发者快速构建高可用、可观测的分布式系统。

## ✨ 核心特性 (Features)

- **🏗 规范架构**: 遵循 Standard Go Project Layout，采用 **DDD (领域驱动设计)** 思想，职责分明 (Interface/Biz/Data)。
- **🔌 依赖注入**: 基于 **Google Wire** 实现的依赖注入，代码解耦，测试友好。
- **🌐 双模通信**: 同时支持 **HTTP (Gin)** 和 **gRPC** 服务，支持优雅退出 (Graceful Shutdown)。
- **🛡 服务治理**:
  - **注册与发现**: 基于 **Etcd** 的自动注册与发现。
  - **负载均衡**: gRPC 客户端侧负载均衡。
  - **熔断限流**: (可扩展) 内置基础拦截器框架。
- **👀 可观测性**:
  - **链路追踪**: 集成 **OpenTelemetry** (Jaeger)，全链路请求追踪。
  - **监控告警**: 集成 **Prometheus** 指标暴露 (Metrics)。
  - **日志系统**: 基于 **Zap** 的高性能结构化日志。

## 🛠 技术栈 (Tech Stack)

| 组件 | 技术选型 | 作用 |
| --- | --- | --- |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) | HTTP 服务与路由 |
| RPC Framework | [gRPC](https://google.golang.org/grpc) | 内部微服务通信 |
| Dependency Injection | [Wire](https://github.com/google/wire) | 依赖注入与代码生成 |
| Config | [Viper](https://github.com/spf13/viper) | 配置文件加载 |
| Registry | [Etcd](https://etcd.io/) | 服务注册与发现 |
| Trace | [OpenTelemetry](https://opentelemetry.io/) | 分布式链路追踪 |
| Metric | [Prometheus](https://prometheus.io/) | 系统监控指标 |
| Logger | [Zap](https://github.com/uber-go/zap) | 高性能日志 |

## 📂 目录结构

```text
Go_skeleton/
├── api/             # Protobuf 定义与生成的代码
├── cmd/             # 程序入口 (main.go, wire.go)
├── configs/         # 配置文件
├── internal/        # 私有业务逻辑
│   ├── biz/         # 业务逻辑层 (Usecase)
│   ├── data/        # 数据访问层 (Repository)
│   ├── service/     # 接口实现层 (Handler)
├── pkg/             # 通用工具库 (Logger, Server封装, Middleware等)
└── docker-compose.yaml # 基础设施启动脚本
```
## 🚀 快速开始 (Quick Start)

1. 启动基础设施
   使用 Docker Compose 一键启动 Etcd, Jaeger, Prometheus：

   docker-compose up -d
2. 下载依赖

   go mod tidy
3. 生成依赖注入代码 (Wire)
   cd cmd/server
   wire
   cd ../..
4. 运行服务
   go run ./cmd/server -f ./configs
5. 验证
   HTTP 接口: POST http://localhost:8080/users
   Metrics: http://localhost:8080/metrics
   Etcd: http://localhost:8081 (EtcdKeeper)
   Jaeger: http://localhost:16686
   Prometheus: http://localhost:9090
   Grafana: http://localhost:3000

## 📝 开发指南
1. 在 api/ 定义 .proto 文件。
2. 在 internal/biz 定义业务实体与接口。
3. 在 internal/data 实现数据接口。
4. 在 internal/service 实现 HTTP/RPC 接口。
5. 在 cmd/server/wire.go 注册新的 Provider。
6. 运行 wire 生成代码并启动。