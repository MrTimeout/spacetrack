package utils

import (
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Pwd = ""
)

func createFile(t *testing.T, dir, fileName string, perm os.FileMode) *os.File {
	var err error

	if strings.TrimSpace(dir) == Pwd {
		if dir, err = os.Getwd(); err != nil {
			t.Fatal(err)
		}
	}

	f, err := os.OpenFile(path.Join(dir, fileName), os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	return f
}

func TestLogger(t *testing.T) {
	var (
		files     = [2]string{"./target_file1", "./target_file2"}
		appenders [2]FileLoggerAppender
	)

	for i := range files {
		appenders[i] = FileLoggerAppender{
			LoggerFileLevel: DebugLevel,
			LoggerFileName:  files[i],
			DateTimeFormat:  RFC3339,
		}
	}

	t.Run("development config with more than one core working properly", func(t *testing.T) {
		t.Cleanup(func() {
			for _, appender := range appenders {
				os.Remove(appender.LoggerFileName)
			}
		})

		want := "This is the message"
		l := Logger{
			Production:    false,
			FileAppenders: appenders[:],
		}

		logger := l.Tee()

		logger.Info(want)

		for _, appender := range appenders {
			b, err := os.ReadFile(appender.LoggerFileName)
			if err != nil {
				t.Fatal(err)
			}

			assert.Contains(t, string(b), want)
		}
	})

	t.Run("production config with more than one core working properly", func(t *testing.T) {
		t.Cleanup(func() {
			for _, appender := range appenders {
				os.Remove(appender.LoggerFileName)
			}
		})

		want := "This is the message"
		l := Logger{
			Production:    true,
			FileAppenders: appenders[:],
		}

		logger := l.Tee()

		logger.Info(want)

		for _, appender := range appenders {
			b, err := os.ReadFile(appender.LoggerFileName)
			if err != nil {
				t.Fatal(err)
			}

			assert.Contains(t, string(b), want)
		}
	})

	t.Run("new logger creation successfully", func(t *testing.T) {
		got := NewLogger(true, DebugLevel, "./some", "./file", "./here")
		want := Logger{
			Production: false,
			FileAppenders: []FileLoggerAppender{
				{LoggerFileLevel: DebugLevel, LoggerFileName: "./some", DateTimeFormat: RFC3339},
				{LoggerFileLevel: DebugLevel, LoggerFileName: "./file", DateTimeFormat: RFC3339},
				{LoggerFileLevel: DebugLevel, LoggerFileName: "./here", DateTimeFormat: RFC3339},
			},
			ConsoleAppender: ConsoleAppender{LoggerFileLevel: DebugLevel, DateTimeFormat: RFC3339},
		}

		assert.Equal(t, want, got)
	})
}

func TestFileLoggerAppender(t *testing.T) {
	t.Run("log file is created and log is wrote correctly", func(t *testing.T) {
		want := "This is the message"
		f := "./target_file"
		t.Cleanup(func() {
			os.Remove(f)
		})

		appender := FileLoggerAppender{
			LoggerFileLevel: DebugLevel,
			LoggerFileName:  f,
			DateTimeFormat:  RFC3339,
		}

		core := appender.core(zap.NewDevelopmentEncoderConfig())

		err := core.Write(zapcore.Entry{
			Level:      zap.InfoLevel,
			Time:       time.Now(),
			LoggerName: "testing",
			Message:    want,
		}, []zapcore.Field{})

		if err != nil {
			t.Fatal(err)
		}

		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, string(b), want)
	})

	t.Run("log file don't have permissions so an empty core is created", func(t *testing.T) {
		f := createFile(t, Pwd, "randomName", 0444)

		appender := FileLoggerAppender{
			LoggerFileLevel: DebugLevel,
			LoggerFileName:  f.Name(),
			DateTimeFormat:  RFC3339,
		}

		assert.IsType(t, zapcore.NewNopCore(), appender.core(zap.NewDevelopmentEncoderConfig()))
	})
}

func TestDateTimeFormat(t *testing.T) {
	var dateTimeFormatCases = []struct {
		input string
		want  DateTimeFormat
	}{
		{
			input: "ansic",
			want:  ANSIC,
		},
		{
			input: "unixDate",
			want:  UnixDate,
		},
		{
			input: "rubyDate",
			want:  RubyDate,
		},
		{
			input: "rfc822",
			want:  RFC822,
		},
		{
			input: "rfc822z",
			want:  RFC822Z,
		},
		{
			input: "rfc850",
			want:  RFC850,
		},
		{
			input: "rfc1123",
			want:  RFC1123,
		},
		{
			input: "rfc3339",
			want:  RFC3339,
		},
		{
			input: "rfc3339nano",
			want:  RFC3339Nano,
		},
		{
			input: "kitchen",
			want:  Kitchen,
		},
		{
			input: "stamp",
			want:  Stamp,
		},
		{
			input: "stampmilli",
			want:  StampMilli,
		},
		{
			input: "stampmicro",
			want:  StampMicro,
		},
		{
			input: "stampnano",
			want:  StampNano,
		},
	}

	t.Run("type of date time format should return string type", func(t *testing.T) {
		want := "string"
		for _, dateTimeFormatCase := range dateTimeFormatCases {
			assert.Equal(t, want, dateTimeFormatCase.want.Type())
		}
	})

	t.Run("set of date time format should set the correct one when input is correct", func(t *testing.T) {
		for _, dateTimeFormatCase := range dateTimeFormatCases {
			var d DateTimeFormat

			assert.Nil(t, d.Set(dateTimeFormatCase.input))
			assert.Equal(t, dateTimeFormatCase.want, d)
		}
	})

	t.Run("set of date time format should return error when incorrect value is passed", func(t *testing.T) {
		var d DateTimeFormat
		assert.ErrorIs(t, ErrDateTimeFormatNotAllowed, d.Set("incorrectDateTimeFormat"))
	})

	t.Run("to zap time encoder should return the zapcore.TimeEncoder correctly", func(t *testing.T) {
		for _, dateTimeFormatCase := range dateTimeFormatCases {
			assert.NotNil(t, dateTimeFormatCase.want.ToZapTimeEncoder())
		}
	})
}

func TestLoggerLevel(t *testing.T) {
	var loggerLevelCases = []struct {
		input string
		want  LoggerLevel
	}{
		{
			input: "debug",
			want:  DebugLevel,
		},
		{
			input: "info",
			want:  InfoLevel,
		},
		{
			input: "warn",
			want:  WarnLevel,
		},
		{
			input: "error",
			want:  ErrorLevel,
		},
		{
			input: "dpanic",
			want:  DPanicLevel,
		},
		{
			input: "panic",
			want:  PanicLevel,
		},
		{
			input: "fatal",
			want:  FatalLevel,
		},
	}

	t.Run("type of logger level should return string type", func(t *testing.T) {
		want := "string"
		for _, loggerLevelCase := range loggerLevelCases {
			assert.Equal(t, want, loggerLevelCase.want.Type())
		}
	})

	t.Run("set of logger level should set the correct one when input is correct", func(t *testing.T) {
		for _, loggerLevelCase := range loggerLevelCases {
			var l LoggerLevel

			assert.Nil(t, l.Set(loggerLevelCase.input))
			assert.Equal(t, loggerLevelCase.want, l)
		}
	})

	t.Run("set of logger level should return error when incorrect value is passed", func(t *testing.T) {
		var l LoggerLevel
		assert.ErrorIs(t, ErrLoggerLevelNotAllowed, l.Set("incorrectLoggerLevel"))
	})

	t.Run("to zap log level should return the zapcore.Level correctly", func(t *testing.T) {
		for _, loggerLevelCase := range loggerLevelCases {
			assert.NotNil(t, loggerLevelCase.want.ToZapLevel())
		}
	})
}

func TestCredentials(t *testing.T) {
	t.Run("spacetrack identity and password are correct secrets", func(t *testing.T) {
		got, err := SpaceTrackAuth{
			Identity: ciphertext,
			Password: ciphertext,
			Secret:   key,
		}.credentials()
		if err != nil {
			t.Fatal(err)
		}
		want := map[string]string{"identity": plaintext, "password": plaintext}

		for key := range want {
			assert.Equal(t, want[key], got[key])
		}
	})

	t.Run("spacetrack with incorrect key size", func(t *testing.T) {
		Configure(Logger{})
		_, gotErr := SpaceTrackAuth{
			Identity: ciphertext,
			Password: ciphertext,
			Secret:   "incorrectkey",
		}.credentials()

		assert.ErrorContains(t, gotErr, "secret")
	})

	t.Run("spacetrack with incorrect identity", func(t *testing.T) {
		_, gotErr := SpaceTrackAuth{
			Identity: "notcorrect",
			Password: ciphertext,
			Secret:   key,
		}.credentials()

		assert.ErrorContains(t, gotErr, "incorrect secret")
	})

	t.Run("spacetrack with incorrect password", func(t *testing.T) {
		_, gotErr := SpaceTrackAuth{
			Identity: ciphertext,
			Password: "notcorrect",
			Secret:   key,
		}.credentials()

		assert.ErrorContains(t, gotErr, "incorrect secret")
	})
}

func TestFormatCookie(t *testing.T) {
	t.Run("cookie is formatted correctly", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{
				Name:  "cookie",
				Value: "cookie-value",
			},
		}

		assert.Equal(t, "cookie=cookie-value", s.FormatCookie())
	})

	t.Run("cookie is nil and returns empty string", func(t *testing.T) {
		s := SpaceTrackAuth{}

		assert.Empty(t, s.FormatCookie())
	})
}

func TestIsAuthNeeded(t *testing.T) {
	t.Run("cookie is nil and returns true", func(t *testing.T) {
		s := SpaceTrackAuth{}

		assert.True(t, s.IsAuthNeeded())
	})

	t.Run("cookie is valid because it has not been expired yet", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{
				Name:    "cookie",
				Value:   "cookie-value",
				MaxAge:  7200,
				Expires: time.Now().Add(-time.Hour),
			},
		}

		assert.False(t, s.IsAuthNeeded())
	})

	t.Run("cookie is invalid because it has been already expired", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{
				Name:    "cookie",
				Value:   "cookie-value",
				MaxAge:  7200,
				Expires: time.Now().Add(-3 * time.Hour),
			},
		}

		assert.True(t, s.IsAuthNeeded())
	})
}

func TestPersistCookies(t *testing.T) {
	t.Run("trying to persist cookie fails because of file with no permissions", func(t *testing.T) {
		Configure(Logger{})
		s := SpaceTrackAuth{}
		f := createTmpWithPerm(t, 0444)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() {
			viper.SetConfigFile(previousViperFile)
		})

		got := s.PersistCookie([]*http.Cookie{{Name: "chocolatechip", Value: "something here"}})

		assert.False(t, got)
	})

	t.Run("persist cookie will not succeed because there is no the right cookie inside the arr", func(t *testing.T) {
		s := SpaceTrackAuth{}

		got := s.PersistCookie([]*http.Cookie{{Name: "nottherightcookie", Value: "some value here"}})

		assert.False(t, got)
	})

	t.Run("persiste cookie correctly", func(t *testing.T) {
		Configure(Logger{})
		fileName := "spacetrack.yml"
		s := SpaceTrackAuth{}
		want := &http.Cookie{Name: "chocolatechip", Value: "some value here"}

		f := createFile(t, Pwd, fileName, 0666)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() { viper.SetConfigFile(previousViperFile) })

		got := s.PersistCookie([]*http.Cookie{want})

		assert.True(t, got)
		assert.Equal(t, want, s.Cookie)
	})
}

func TestEncode(t *testing.T) {
	t.Run("encode successful", func(t *testing.T) {
		s := SpaceTrackAuth{
			Identity: "xxxx",
			Password: "yyyy",
		}
		want := "identity=xxxx&password=yyyy"

		got, err := s.Encode()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, want, got)
	})

	t.Run("encode unsuccesful because of credentials error", func(t *testing.T) {
		_, gotErr := SpaceTrackAuth{
			Identity: "xxx",
			Password: "yyy",
			Secret:   "keysizerror",
		}.Encode()

		assert.ErrorContains(t, gotErr, "incorrect secret")
	})
}
