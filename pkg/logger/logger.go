package logger

import (
	"Go_skeleton/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局变量
// 这里的 *zap.Logger 是 Zap 的核心对象
var Log *zap.Logger

// InitLogger 初始化日志
// 传入之前加载好的配置 config.LogConfig
func InitLogger(cfg config.LogConfig) {

	// 设置日志切割规则
	// lumberjack 会自动管理日志文件 也就是切割日志文件。
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,   // 日志文件路径
		MaxSize:    cfg.MaxSize,    // 每个文件最大尺寸 (MB)
		MaxBackups: cfg.MaxBackups, // 保留备份个数
		MaxAge:     cfg.MaxAge,     // 保留天数
		Compress:   true,           // 是否压缩备份文件
	})

	// 设置日志编码格式
	// 这里配置成 JSON 格式，
	encoderConfig := zap.NewProductionEncoderConfig()
	// 把时间修改成便于阅读的格式：2023-10-01 12:00:00
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 把日志级别大写，比如 INFO, ERROR
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 解析日志级别
	// 配置文件里写的是 string ("debug"), 转成 zap 能识别的 AtomicLevel
	var l = new(zapcore.Level)
	if err := l.UnmarshalText([]byte(cfg.Level)); err != nil {
		// 如果填写级别错误了，默认Info 级别
		*l = zapcore.InfoLevel
	}

	// 组装核心 Core
	core := zapcore.NewCore(encoder, writeSyncer, l)
	// 创建 Logger 对象
	Log = zap.New(core, zap.AddCaller())

	// 替换全局的 logger，用 zap.L()
	zap.ReplaceGlobals(Log)
}
