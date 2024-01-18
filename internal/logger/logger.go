package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	AppLogger *zap.Logger
	ZMQLogger *zap.Logger
}

var Loggers *Logger

func NewLogger(filename string) *zap.Logger {
	// 配置 zap
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置日志文件滚动
	logFile := &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    500,      // 每个日志文件的最大尺寸，单位MB
		MaxBackups: 3,        // 保留的旧日志文件的最大个数
		MaxAge:     3,        // 保留的旧日志文件的最大天数
		Compress:   true,     // 是否压缩旧日志文件
	}
	// 设置日志输出为文件
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(logFile),
		config.Level,
	)

	// 创建 logger
	return zap.New(core)
}

func init() {
	Loggers = &Logger{}
	Loggers.AppLogger = NewLogger("logs/app.log")
	Loggers.ZMQLogger = NewLogger("logs/zmq.log")
}

// ----------------------------------------------------------------
// 项目日志
func Debug(msg string) {
	Loggers.AppLogger.Debug(msg)
}

func Info(msg string) {
	Loggers.AppLogger.Info(msg)
}

func Error(msg string) {
	Loggers.AppLogger.Error(msg)
}

func Fatal(msg string) {
	Loggers.AppLogger.Fatal(msg)
}

// ----------------------------------------------------------------
// ZMQ日志
func ZMQDebug(msg string) {
	Loggers.ZMQLogger.Debug(msg)
}

func ZMQInfo(msg string) {
	Loggers.ZMQLogger.Info(msg)
}

func ZMQError(msg string) {
	Loggers.ZMQLogger.Error(msg)
}

func ZMQFatal(msg string) {
	Loggers.ZMQLogger.Fatal(msg)
}
