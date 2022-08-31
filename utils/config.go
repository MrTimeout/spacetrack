package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// ErrDateTimeFormatNotAllowed is used to indicate to the user that the format which is trying to use is not allowed
	// for zap logger. Possible values are in this file and there are quite of them.
	ErrDateTimeFormatNotAllowed = errors.New("date time format not allowed for zap logger")
	// ErrLoggerLevelNotAllowed is used to indicate that the logger level the user has specified is not allowed.
	ErrLoggerLevelNotAllowed = errors.New("logger level not allowed")
	// ErrIncorrectSecret is thrown when the secret to decrypt credentials is incorrect
	ErrIncorrectSecret = errors.New("incorrect secret")
)

// Config is the main structure which we are going to use to store the configuration of the application.
// Here we have the logger configuration.
type Config struct {
	Auth       SpaceTrackAuth `json:"auth" yaml:"auth" mapstructure:"auth"`
	WorkDir    string         `json:"work_dir" yaml:"work_dir" mapstructure:"work_dir"`
	Interval   string         `json:"interval" yaml:"interval" mapstructure:"interval"`
	SecretFile string         `json:"secret_file" yaml:"secret_file" mapstructure:"secret_file"`
	Logger     Logger         `json:"logger" yaml:"logger" mapstructure:"logger"`
}

// SpaceTrackAuth
type SpaceTrackAuth struct {
	Identity string       `json:"identity" yaml:"identity" mapstructure:"identity"`
	Password string       `json:"password" yaml:"password" mapstructure:"password"`
	Cookie   *http.Cookie `json:"cookie" yaml:"cookie" mapstructure:"cookie"`
	Secret   string       `yaml:"-"`
}

// Encode will encode the credentials to pass them along the http request
func (sta SpaceTrackAuth) Encode() (string, error) {
	var another url.Values = url.Values{}

	credentials, err := sta.credentials()
	if err != nil {
		return "", err
	}

	for k, v := range credentials {
		another[k] = []string{v}
	}

	return another.Encode(), nil
}

func (sta SpaceTrackAuth) credentials() (map[string]string, error) {
	cred := map[string]string{
		"identity": sta.Identity,
		"password": sta.Password,
	}

	if sta.Secret != "" {
		if err := sta.decrypt("identity", sta.Identity, cred); err != nil {
			return nil, err
		}
		if err := sta.decrypt("password", sta.Password, cred); err != nil {
			return nil, err
		}
	}

	return cred, nil
}

func (sta SpaceTrackAuth) decrypt(name, value string, credentials map[string]string) error {
	if decrypted, err := Decrypt([]byte(value), []byte(sta.Secret)); err != nil {
		Warn("incorrect secret", zap.Error(err))
		return ErrIncorrectSecret
	} else {
		credentials[name] = string(decrypted)
	}
	return nil
}

// PersistCookie will persist cookie chocolatechip if found
func (sta *SpaceTrackAuth) PersistCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == "chocolatechip" {
			viper.Set("auth.cookie", cookie)
			sta.Cookie = cookie
			if err := viper.WriteConfig(); err != nil {
				Warn("trying to write the cookies to the config file", zap.Any("cookie", cookie), zap.Error(err))
				return false
			} else {
				Info("We have successfully update the cookie value in the configuration file")
				return true
			}
		}
	}
	return false
}

// FormatCookie will format the cookie as the standard allows us
func (sta SpaceTrackAuth) FormatCookie() string {
	if sta.Cookie != nil {
		return fmt.Sprintf("%s=%s", sta.Cookie.Name, sta.Cookie.Value)
	}
	return ""
}

// IsAuthNeeded returns if the cookie has expired or null being true the result.
func (sta SpaceTrackAuth) IsAuthNeeded() bool {
	return sta.Cookie == nil || sta.Cookie.Expires.Add(
		time.Duration(sta.Cookie.MaxAge)*time.Second).Before(time.Now())
}

// Logger is where all zap logger stuff will go
type Logger struct {
	Production      bool                 `json:"prod" yaml:"prod" mapstructure:"prod"`
	FileAppenders   []FileLoggerAppender `json:"file_appenders" yaml:"file_appenders" mapstructure:"file_appenders"`
	ConsoleAppender ConsoleAppender      `json:"console_appender" yaml:"console_appender" mapstructure:"console_appender"`
}

// NewLogger returns a new Logger with a logger level and some files
func NewLogger(console bool, loggerLevel LoggerLevel, loggerFileNames ...string) Logger {
	var (
		appenders       = make([]FileLoggerAppender, len(loggerFileNames))
		consoleAppender ConsoleAppender
	)

	for i := range loggerFileNames {
		appenders[i] = NewFileLoggerAppender(loggerLevel, loggerFileNames[i], RFC3339)
	}

	if console {
		consoleAppender = NewConsoleAppender(loggerLevel)
	}

	return Logger{
		Production:      false,
		FileAppenders:   appenders,
		ConsoleAppender: consoleAppender,
	}
}

// Tee create core loggers to log into them
func (l Logger) Tee() *zap.Logger {
	var cfg zapcore.EncoderConfig

	if l.Production {
		cfg = zap.NewProductionEncoderConfig()
	} else {
		cfg = zap.NewDevelopmentEncoderConfig()
	}

	cores := make([]zapcore.Core, len(l.FileAppenders))
	for i := range l.FileAppenders {
		cores[i] = l.FileAppenders[i].core(cfg)
	}

	cores = append(cores, l.ConsoleAppender.core(cfg))

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// Appender describes a standard appender to the zap logger.
type Appender interface {
	// ... implements all the logic which can help us to create the zap logger.
	core(zapcore.EncoderConfig) zapcore.Core
}

// ConsoleAppender is the struct that allows to add a console appender to the zap logger
type ConsoleAppender struct {
	LoggerFileLevel LoggerLevel    `json:"level" yaml:"level" mapstructure:"level"`
	DateTimeFormat  DateTimeFormat `json:"date_format" yaml:"date_format" mapstructure:"date_format"`
}

// NewConsoleAppender returns a ConsoleAppender with logger level specified
func NewConsoleAppender(loggerLevel LoggerLevel) ConsoleAppender {
	return ConsoleAppender{
		LoggerFileLevel: loggerLevel,
		DateTimeFormat:  RFC3339,
	}
}

func (ca ConsoleAppender) core(config zapcore.EncoderConfig) zapcore.Core {
	config.EncodeTime = ca.DateTimeFormat.ToZapTimeEncoder()
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), ca.LoggerFileLevel.ToZapLevel())
}

// FileLoggerAppender is the struct that allows to add a file appender
// to the zap logger
type FileLoggerAppender struct {
	LoggerFileLevel LoggerLevel    `json:"level" yaml:"level" mapstructure:"level"`
	LoggerFileName  string         `json:"file" yaml:"file" mapstructure:"file"`
	DateTimeFormat  DateTimeFormat `json:"date_format" yaml:"date_format" mapstructure:"date_format"`
}

// NewFileLoggerAppender returns a FileLoggerAppender with values passed as parameters
func NewFileLoggerAppender(loggerLevel LoggerLevel, fileName string, dateTimeFormat DateTimeFormat) FileLoggerAppender {
	return FileLoggerAppender{
		LoggerFileLevel: loggerLevel,
		LoggerFileName:  fileName,
		DateTimeFormat:  dateTimeFormat,
	}
}

func (fla FileLoggerAppender) core(config zapcore.EncoderConfig) zapcore.Core {
	if logfile, err := os.OpenFile(fla.LoggerFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660); err == nil {
		config.EncodeTime = fla.DateTimeFormat.ToZapTimeEncoder()
		return zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(logfile), fla.LoggerFileLevel.ToZapLevel())
	}
	return zapcore.NewNopCore()
}

// DateTimeFormat is just a string type, that contains all the date time formats allowed by zap library.
type DateTimeFormat string

const (
	ANSIC       = "ansic"
	UnixDate    = "unixDate"
	RubyDate    = "rubyDate"
	RFC822      = "rfc822"
	RFC822Z     = "rfc822z"
	RFC850      = "rfc850"
	RFC1123     = "rfc1123"
	RFC1123Z    = "rfc1123z"
	RFC3339     = "rfc3339"
	RFC3339Nano = "rfc3339nano"
	Kitchen     = "kitchen"
	Stamp       = "stamp"
	StampMilli  = "stampMilli"
	StampMicro  = "stampMicro"
	StampNano   = "stampNano"
)

// Type give us the type of DateTimeFormat. It is useful when dealing with flags in go.
func (d DateTimeFormat) Type() string {
	return "string"
}

// Set tries to set the value checking if the input is correct and returning error ErrDateTimeFormatNotAllowed otherwise.
func (d *DateTimeFormat) Set(input string) error {
	if !d.unmarshal(strings.ToLower(input)) {
		return ErrDateTimeFormatNotAllowed
	}
	return nil
}

// nolint:cyclop
func (d *DateTimeFormat) unmarshal(input string) bool {
	switch input {
	case "ansic":
		*d = ANSIC
	case "unixdate":
		*d = UnixDate
	case "rubydate":
		*d = RubyDate
	case "rfc822":
		*d = RFC822
	case "rfc822z":
		*d = RFC822Z
	case "rfc850":
		*d = RFC850
	case "rfc1123":
		*d = RFC1123
	case "rfc3339":
		*d = RFC3339
	case "rfc3339nano":
		*d = RFC3339Nano
	case "kitchen":
		*d = Kitchen
	case "stamp":
		*d = Stamp
	case "stampmilli":
		*d = StampMilli
	case "stampmicro":
		*d = StampMicro
	case "stampnano":
		*d = StampNano
	default:
		return false
	}

	return true
}

// nolint:cyclop
// String is the string representation of the date time format.
func (d *DateTimeFormat) String() string {
	var result = ""

	switch *d {
	case ANSIC:
		result = time.ANSIC
	case UnixDate:
		result = time.UnixDate
	case RubyDate:
		result = time.RubyDate
	case RFC822:
		result = time.RFC822
	case RFC822Z:
		result = time.RFC822Z
	case RFC850:
		result = time.RFC850
	case RFC1123:
		result = time.RFC1123
	case RFC3339:
		result = time.RFC3339
	case RFC3339Nano:
		result = time.RFC3339Nano
	case Kitchen:
		result = time.Kitchen
	case Stamp:
		result = time.Stamp
	case StampMilli:
		result = time.StampMilli
	case StampMicro:
		result = time.StampMicro
	case StampNano:
		result = time.StampNano
	}

	return result
}

// ToZapTimeEncoder returns the zapcore.TimeEncoder so we can use it to set it into the zap.logger configuration
func (d *DateTimeFormat) ToZapTimeEncoder() zapcore.TimeEncoder {
	return zapcore.TimeEncoderOfLayout(d.String())
}

// LoggerLevel is just a wrapper of the zapcore.Level type of zap library which we
// are going to use with cobra and viper
type LoggerLevel string

const (
	// DebugLevel just for development.
	DebugLevel = "debug"
	// InfoLevel is the default logging priority.
	InfoLevel = "info"
	// WarnLevel is more important logs, but don't require human review.
	WarnLevel = "warn"
	// ErrorLevel logs of high-priority.
	ErrorLevel = "error"
	// DPanicLevel logs are particularly important errors.
	DPanicLevel = "dpanic"
	// PanicLevel logs a message, then panics.
	PanicLevel = "panic"
	// FatalLevel logs a message, then calls os.Exit(1)
	FatalLevel = "fatal"
)

// Type returns the type of the LoggerLevel type
func (l *LoggerLevel) Type() string {
	return "string"
}

// Set tries to set the LoggerLevel returning error if the input is incorrect
func (l *LoggerLevel) Set(input string) error {
	if !l.unmarshal(strings.ToLower(input)) {
		return ErrLoggerLevelNotAllowed
	}
	return nil
}

func (l *LoggerLevel) unmarshal(input string) bool {
	switch input {
	case "debug":
		*l = DebugLevel
	case "info":
		*l = InfoLevel
	case "warn":
		*l = WarnLevel
	case "error":
		*l = ErrorLevel
	case "dpanic":
		*l = DPanicLevel
	case "panic":
		*l = PanicLevel
	case "fatal":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

// String is the string representation of the LoggerLevel
func (l *LoggerLevel) String() string {
	var result = ""

	switch *l {
	case DebugLevel:
		result = "debug"
	case InfoLevel:
		result = "info"
	case WarnLevel:
		result = "warn"
	case ErrorLevel:
		result = "error"
	case DPanicLevel:
		result = "dpanic"
	case PanicLevel:
		result = "panic"
	case FatalLevel:
		result = "fatal"
	}

	return result
}

// ToZapLevel is used to parse our custom LoggerLevel to the zapcore.Level
// so we can use it to our zap.Logger configuration
func (l *LoggerLevel) ToZapLevel() zapcore.Level {
	z, _ := zapcore.ParseLevel(l.String()) //nolint:errcheck
	return z
}
