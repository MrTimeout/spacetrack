package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func fileContains(t *testing.T, fileName string, contains string) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, string(b), contains)
}

func TestLoggerFuncs(t *testing.T) {
	f := createFile(t, "", "spacetrack.log", 0666)
	Configure(Logger{
		Production: false,
		FileAppenders: []FileLoggerAppender{
			{
				LoggerFileLevel: DebugLevel,
				LoggerFileName:  f.Name(),
				DateTimeFormat:  RFC3339,
			},
		},
	})

	for _, each := range []struct {
		description, msg string
		loggerLevel      LoggerLevel
		logLevelFunc     func(string, ...zap.Field)
	}{
		{
			description:  "Log into file using debug logger level",
			msg:          "debug message here",
			loggerLevel:  DebugLevel,
			logLevelFunc: Debug,
		},
		{
			description:  "Log into file using info logger level",
			msg:          "info message here",
			loggerLevel:  InfoLevel,
			logLevelFunc: Info,
		},
		{
			description:  "Log into file using warn logger level",
			msg:          "warn message here",
			loggerLevel:  WarnLevel,
			logLevelFunc: Warn,
		},
		{
			description:  "Log into file using error logger level",
			msg:          "error message here",
			loggerLevel:  ErrorLevel,
			logLevelFunc: Error,
		},
		{
			description:  "Log into file using dpanic logger level",
			msg:          "dpanic message here",
			loggerLevel:  DPanicLevel,
			logLevelFunc: DPanic,
		},
		{
			description:  "Log into file using panic logger level",
			msg:          "panic message here",
			loggerLevel:  PanicLevel,
			logLevelFunc: Panic,
		},
	} {
		t.Run(each.description, func(t *testing.T) {
			if each.loggerLevel == PanicLevel {
				assert.Panics(t, func() { each.logLevelFunc(each.msg) })
			} else {
				each.logLevelFunc(each.msg)
			}

			fileContains(t, f.Name(), each.msg)
		})
	}
}
