package logger

import (
	"context"

	zlog "github.com/rs/zerolog/log"
)

// GetCorrelationIDLoggerCtxFromRequest gets new context logger with correlation ID
func GetCorrelationIDLoggerCtx(ctx context.Context, cid string) context.Context {
	cidlog := zlog.With().Str(LOG_KEY_CORRELATION_ID, cid).Logger()
	return context.WithValue(ctx, LoggerCtxKey, cidlog)
}
