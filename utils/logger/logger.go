package logger

import (
	"context"
	"fmt"
	"golang-boiler-plate/utils/types"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type LoggerWrapper interface {
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	InboundRequest(url string) *loggerWrapper
	OutboundRequest(url string) *loggerWrapper
	AdditionalInfo(info types.Obj) *loggerWrapper
	StackTrace(err error) *loggerWrapper
}

type LogCtxKey int

const (
	LoggerCtxKey LogCtxKey = iota
)

var (
	INBOUND_REQUEST_PREFIX  = "INBOUND"
	OUTBOUND_REQUEST_PREFIX = "OUTBOUND"
	EXCLUDE_PATH            = [...]string{
		"/status",
		"/ping",
		"/health",
	}
)

var (
	LOG_KEY_CORRELATION_ID  = "correlationId"
	LOG_KEY_RESPONSE_TIME   = "responseTime"
	LOG_KEY_STACK_TRACE     = "stackTrace"
	LOG_KEY_ADDITIONAL_INFO = "additionalInfo"
	LOG_KEY_STATUS_CODE     = "status"
)

type loggerWrapper struct {
	logger           zerolog.Logger
	logbase          *zerolog.Event
	startTime        time.Time
	messagePrefix    string
	isInExcludedPath bool
	stackTrace       error
	additionalInfo   types.Obj
}

func Ctx(ctx context.Context) LoggerWrapper {
	logger := zerolog.Logger{}
	ctxVal := ctx.Value(LoggerCtxKey)
	if ctxVal != nil {
		logger = ctx.Value(LoggerCtxKey).(zerolog.Logger)
	}
	return &loggerWrapper{
		logger:           logger,
		startTime:        time.Now(),
		additionalInfo:   nil,
		messagePrefix:    "",
		isInExcludedPath: false,
		stackTrace:       nil,
	}
}

func ConfigureLogger(logFormat bool, debugMode bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if logFormat {
		configureProdLogger()
	} else {
		configureLocalLogger()
	}
}

func configureProdLogger() {
	zerolog.LevelInfoValue = "INFO"
	zerolog.LevelDebugValue = "DEBUG"
	zerolog.LevelErrorValue = "ERROR"
	zerolog.LevelFatalValue = "FATAL"
	zerolog.LevelWarnValue = "WARN"
	zerolog.LevelPanicValue = "PANIC"
	zerolog.ErrorStackFieldName = "stackTrace"
}

func configureLocalLogger() {
	output := zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatTimestamp: func(i interface{}) string {
			parse, _ := time.Parse(time.RFC3339, i.(string))
			return parse.Format("2006-01-02 15:04:05")
		},
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf(" %-6s ", i))
		},
	}

	zlog.Logger = zerolog.New(output).With().
		Timestamp().CallerWithSkipFrameCount(4).Logger()
}

func ExcludedPath(url string) bool {
	for _, path := range EXCLUDE_PATH {
		if path == url {
			return true
		}
	}
	return false
}

func Debug(format string, args ...interface{}) {
	zlog.Debug().Msgf(format, args...)
}

// log info level without context
func Info(format string, args ...interface{}) {
	zlog.Info().Msgf(format, args...)
}

// log error level without context
func Error(format string, args ...interface{}) {
	zlog.Error().Msgf(format, args...)
}

// log fatal level without context
func Fatal(format string, args ...interface{}) {
	zlog.Fatal().Msgf(format, args...)
}

func (l *loggerWrapper) logFormatter(format string, args ...interface{}) {
	if l.messagePrefix != "" {
		format = l.messagePrefix + " " + format
	}

	isInboundRequest := l.messagePrefix == INBOUND_REQUEST_PREFIX
	isOutboundRequest := l.messagePrefix == OUTBOUND_REQUEST_PREFIX
	isHttpLog := strings.Contains(format, "[httplog]")

	if (isOutboundRequest && !isHttpLog) || (isInboundRequest && isHttpLog) {
		l.AdditionalInfo(types.Obj{
			LOG_KEY_RESPONSE_TIME: time.Since(l.startTime).Milliseconds(),
		})
	}

	if l.stackTrace != nil {
		l.AdditionalInfo(types.Obj{
			LOG_KEY_STACK_TRACE: fmt.Sprintf("%+v", l.stackTrace),
		})
	}

	if len(l.additionalInfo) > 0 {
		l.logbase.Interface(LOG_KEY_ADDITIONAL_INFO, l.additionalInfo)
	}

	if !l.isInExcludedPath {
		l.logbase.Msgf(format, args...)
	}

	l.additionalInfo = nil
	l.stackTrace = nil
	l.messagePrefix = ""
}

func (l *loggerWrapper) Info(format string, args ...interface{}) {
	l.logbase = l.logger.Info()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Debug(format string, args ...interface{}) {
	l.logbase = l.logger.Debug()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Warn(format string, args ...interface{}) {
	l.logbase = l.logger.Warn()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Error(format string, args ...interface{}) {
	l.logbase = l.logger.Error()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Fatal(format string, args ...interface{}) {
	l.logbase = l.logger.Fatal()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) InboundRequest(url string) *loggerWrapper {
	l.messagePrefix = INBOUND_REQUEST_PREFIX
	l.isInExcludedPath = isUrlIsInExcludedPath(url)
	return l
}

func (l *loggerWrapper) OutboundRequest(url string) *loggerWrapper {
	l.messagePrefix = OUTBOUND_REQUEST_PREFIX
	l.isInExcludedPath = isUrlIsInExcludedPath(url)
	return l
}

func (l *loggerWrapper) AdditionalInfo(info types.Obj) *loggerWrapper {
	if l.additionalInfo == nil {
		l.additionalInfo = info
	} else {
		for k, v := range info {
			l.additionalInfo[k] = v
		}
	}

	return l
}

func (l *loggerWrapper) StackTrace(err error) *loggerWrapper {
	l.stackTrace = err
	return l
}

func isUrlIsInExcludedPath(url string) bool {
	for _, path := range EXCLUDE_PATH {
		if strings.Contains(url, path) {
			return true
		}
	}

	return false
}
