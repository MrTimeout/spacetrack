package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

var loggerLevelTests = []struct {
	loggerLevel LoggerLevel
	zapLevel    zapcore.Level
}{
	{
		loggerLevel: LoggerLevel{value: "debug"},
		zapLevel:    zapcore.DebugLevel,
	},
	{
		loggerLevel: LoggerLevel{value: "info"},
		zapLevel:    zapcore.InfoLevel,
	},
	{
		loggerLevel: LoggerLevel{value: "warn"},
		zapLevel:    zapcore.WarnLevel,
	},
	{
		loggerLevel: LoggerLevel{value: "error"},
		zapLevel:    zapcore.ErrorLevel,
	},
	{
		loggerLevel: LoggerLevel{value: "dpanic"},
		zapLevel:    zapcore.DPanicLevel,
	},
	{
		loggerLevel: LoggerLevel{value: "panic"},
		zapLevel:    zapcore.PanicLevel,
	},
	{
		loggerLevel: LoggerLevel{value: "fatal"},
		zapLevel:    zapcore.FatalLevel,
	},
}

func createTmpFile(t *testing.T) string {
	f, err := os.CreateTemp(".", "zaploggerfile")
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	return f.Name()
}

func TestInitLogger(t *testing.T) {
	t.Run("zap logger when console is on", func(t *testing.T) {
		InitLogger("", LoggerLevel{value: "info"}, true)
		assert.NotNil(t, Logger)
	})

	t.Run("zap logger when console is off", func(t *testing.T) {
		InitLogger("", LoggerLevel{value: "info"}, false)
		assert.NotNil(t, Logger)
	})

	t.Run("zap logger when file exists", func(t *testing.T) {
		InitLogger(createTmpFile(t), LoggerLevel{value: "info"}, false)
		assert.NotNil(t, Logger)
	})
}

func TestLoggerLevel(t *testing.T) {
	t.Run("map logger level to zap", func(t *testing.T) {
		for _, l := range loggerLevelTests {
			assert.Equal(t, l.zapLevel, l.loggerLevel.Zap())
		}
	})

	t.Run("string representation of logger level same as zap", func(t *testing.T) {
		for _, l := range loggerLevelTests {
			assert.Equal(t, l.zapLevel.String(), l.loggerLevel.String())
		}
	})

	t.Run("string representation of logger level is parsed correctly", func(t *testing.T) {
		var got LoggerLevel
		for _, l := range loggerLevelTests {
			err := got.Set(l.loggerLevel.value)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, l.zapLevel, got.Zap())
		}
	})

	t.Run("error when trying to set invalid logger level", func(t *testing.T) {
		var l LoggerLevel
		got := l.Set("invalid field")
		want := "unrecognized level: \"invalid field\""

		assert.EqualError(t, got, want)
	})

	t.Run("type logger level is string", func(t *testing.T) {
		for _, l := range loggerLevelTests {
			assert.Equal(t, "string", l.loggerLevel.Type())
		}
	})
}
