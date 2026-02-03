package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 我们定义两个最常用的指标
// 1. 请求计数器 (Counter)
// 记录一共收到了多少个请求，是个只增不减的数字
var RequestCounter = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total", // 指标名字，Prometheus 里搜这个
		Help: "HTTP 请求总数",
	},
	// 标签：我们可以按 method (GET/POST) 和 path (/user/login) 来分类统计
	[]string{"method", "path"},
)

// 2. 耗时直方图 (Histogram)
// 记录请求处理花了多久
var RequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP 请求耗时分布",
		// 桶 (Buckets)：把耗时分成几档
		// <0.1s, <0.3s, <1.2s, <5s, >5s
		Buckets: []float64{0.1, 0.3, 1.2, 5.0},
	},
	[]string{"method", "path"},
)

// 这里不需要写 Init 函数。
// promauto 这个包很智能，只要代码加载了，它就会自动注册这些指标。
