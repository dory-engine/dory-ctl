package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// error logger
var ZapLog_V1 *zap.Logger

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, //这里可以指定颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.InfoLevel)
	config := zap.Config{
		Level:       atom,  // 日志级别
		Development: false, // 开发模式，堆栈跟踪
		//Encoding:         "json",                                              // 输出格式 console 或 json
		Encoding:         "console",                                            // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                                        // 编码器配置
		InitialFields:    map[string]interface{}{"serviceName": "wisdom_park"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"},                                   // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder //这里可以指定颜色
	// 构建日志
	var err error
	ZapLog_V1, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}

}

func main() {
	ZapLog_V1.Warn("log 初始化成功")
}
