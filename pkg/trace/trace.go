package trace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger" // 导出到 Jaeger
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer 初始化全局 Tracer
// serviceName: 你的服务名叫啥（比如 user-service）
// jaegerEndpoint: Jaeger 的收集地址（比如 http://localhost:14268/api/traces）
func InitTracer(serviceName string, jaegerEndpoint string) (func(context.Context) error, error) {

	// 创建导出器 (Exporter)
	// 负责把生成的 Trace 数据发给 Jaeger
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, fmt.Errorf("创建 Jaeger Exporter 失败: %w", err)
	}

	// 创建资源 (Resource)
	// 给服务贴标签，告诉 Jaeger 这些数据是谁产生的
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// 服务名称必须填写。否则在 Jaeger 列表里找不到
			semconv.ServiceName(serviceName),
			// 可以加版本号、环境等
			attribute.String("environment", "dev"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("创建资源标签失败: %w", err)
	}

	// 创建 TracerProvider
	// 核心管理所有的 Tracer
	tp := sdktrace.NewTracerProvider(
		// 采样率：AlwaysSample 意思是 100% 记录。
		// 生产环境流量大时，通常只采 10% (sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1)))
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		// 绑定身份标签
		sdktrace.WithResource(res),
	)

	// 设置全局 Provider
	// 这样以后在任何地方调用 otel.Tracer("xxx") 都能拿到了
	otel.SetTracerProvider(tp)

	// 设置传播器
	// 它的作用是把 TraceID 塞进 HTTP Header 或者 gRPC Metadata 里，
	// 传给下一个服务。没有它，链路就断了。
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // 支持 W3C 标准
		propagation.Baggage{},      // 支持携带自定义数据
	))

	// 返回一个关闭函数，Server 退出的时候要调用，把还没发出去的数据发完
	return tp.Shutdown, nil
}
