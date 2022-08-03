package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	InitLogger("", LoggerLevel{value: "info"}, true)
}

func InitLogger(logFileName string, logLevel LoggerLevel, console bool) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		createCoreForFile(config, logFileName, logLevel),
		createCoreForConsole(config, logLevel, console),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func createCoreForFile(config zapcore.EncoderConfig, logFileName string, logLevel LoggerLevel) zapcore.Core {
	if logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		fileEncoder := zapcore.NewJSONEncoder(config)
		writer := zapcore.AddSync(logFile)
		return zapcore.NewCore(fileEncoder, writer, logLevel.Zap())
	}

	return zapcore.NewNopCore()
}

func createCoreForConsole(config zapcore.EncoderConfig, logLevel LoggerLevel, console bool) zapcore.Core {
	if console {
		consoleEncoder := zapcore.NewConsoleEncoder(config)
		return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel.Zap())
	}
	return zapcore.NewNopCore()
}

type LoggerLevel struct {
	value string
}

func (l *LoggerLevel) Zap() zapcore.Level {
	level, _ := zapcore.ParseLevel(l.value)
	return level
}

func (l *LoggerLevel) String() string {
	return l.Zap().String()
}

func (l *LoggerLevel) Set(level string) error {
	zapLevel, err := zapcore.ParseLevel(level)
	if err == nil {
		l.value = zapLevel.String()
	}
	return err
}

func (l *LoggerLevel) Type() string {
	return "string"
}
