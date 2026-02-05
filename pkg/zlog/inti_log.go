package zlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

func InitLog(env, level string) {
	//创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 级别大写显示：
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径显示
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	var encoder zapcore.Encoder
	if env == "dev" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //开发模式就彩色显示
		//文本编码器
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		//JSON编码器
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	//配置日志级别
	logLevel := zapcore.DebugLevel
	if level == "info" {
		logLevel = zapcore.InfoLevel
	} else if level == "warn" {
		logLevel = zapcore.WarnLevel
	} else if level == "error" {
		logLevel = zapcore.ErrorLevel
	}

	//配置输出目标
	writeSyncer := zapcore.AddSync(os.Stdout)

	//构建 zap core
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	//构建 logger
	logger = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)

	logger.Info("日志初始化成功")
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Panicf(format string, v ...interface{}) {
	logger.Panic(fmt.Sprintf(format, v...))
}
