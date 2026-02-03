package logger

import (
	"Go_skeleton/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局变量，以后可以在别的代码里直接 logger.Log 就能用
// 这里的 *zap.Logger 是 Zap 的核心对象
var Log *zap.Logger

// InitLogger 初始化日志
// 传入之前加载好的配置 config.LogConfig  也就是"Go_skeleton/pkg/config" 的导入的包
func InitLogger(cfg config.LogConfig) {

	// 1. 设置日志切割规则
	// lumberjack 会自动管理日志文件 也就是切割日志文件。
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,   // 日志文件路径
		MaxSize:    cfg.MaxSize,    // 每个文件最大尺寸 (MB)
		MaxBackups: cfg.MaxBackups, // 保留备份个数
		MaxAge:     cfg.MaxAge,     // 保留天数
		Compress:   true,           // 是否压缩备份文件
	})

	// 2. 设置日志编码格式
	// 这里配置成 JSON 格式，
	encoderConfig := zap.NewProductionEncoderConfig()
	// 把时间修改成便于阅读的格式：2023-10-01 12:00:00
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 把日志级别大写，比如 INFO, ERROR
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 3. 解析日志级别
	// 配置文件里写的是 string ("debug"), 我们要把它转成 zap 能识别的 AtomicLevel
	var l = new(zapcore.Level)
	if err := l.UnmarshalText([]byte(cfg.Level)); err != nil {
		// 如果填写级别错误了，我们这里就默认Info 级别
		*l = zapcore.InfoLevel
	}

	// 4. 组装核心 Core
	// Core 是连接 编码器、写入器、日志级别的桥梁
	core := zapcore.NewCore(encoder, writeSyncer, l)

	// 5. 创建 Logger 对象
	// zap.AddCaller() 的作用是，打印日志时，会显示是哪行代码打印的（文件名:行号）
	Log = zap.New(core, zap.AddCaller())

	// 替换全局的 logger，习惯用 zap.L() (如果你和我一样习惯用zap.L(),那么这一步是必须的)
	zap.ReplaceGlobals(Log)
}

// 补充：为了方便，其实可以再封装一些 Info, Error 方法，
// 这里我选择直接用 zap.L().Info() 也是最标准的写法。
