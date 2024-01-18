package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func init() {
	// 配置 zap
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置日志文件滚动
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    500,            // 每个日志文件的最大尺寸，单位MB
		MaxBackups: 3,              // 保留的旧日志文件的最大个数
		MaxAge:     3,              // 保留的旧日志文件的最大天数
		Compress:   true,           // 是否压缩旧日志文件
	}
	// 设置日志输出为文件
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(logFile),
		config.Level,
	)

	// 创建 logger
	Logger = zap.New(core)
}

func Debug(msg string) {
	Logger.Info(msg)
}

func Info(msg string) {
	Logger.Info(msg)
}

func Error(msg string) {
	Logger.Error(msg)
}

func Fatal(msg string) {
	Logger.Fatal(msg)
}
