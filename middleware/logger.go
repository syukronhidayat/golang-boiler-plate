package middleware

import (
	"context"
	"fmt"
	"golang-boiler-plate/utils/logger"
	"golang-boiler-plate/utils/types"
	"net/http"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		msg := fmt.Sprintf("[%s] %s", r.Method, url)
		lgr := logger.Ctx(r.Context())

		lgr.InboundRequest(url).Info(msg)

		customResponseWriter := NewLoggingResponseWriter(w)

		defer func(ctx context.Context, lrw *LoggingResponseWriter) {
			select {
			case <-ctx.Done():
				lgr.OutboundRequest(url).AdditionalInfo(types.Obj{
					logger.LOG_KEY_STATUS_CODE: http.StatusInternalServerError,
				}).Warn(msg)
			default:
				statusCode := lrw.statusCode
				subLog := lgr.OutboundRequest(url).AdditionalInfo(types.Obj{
					logger.LOG_KEY_STATUS_CODE: statusCode,
				})

				if statusCode < 400 {
					subLog.Info(msg)
				} else {
					subLog.Error(msg)
				}
			}
		}(r.Context(), customResponseWriter)

		next.ServeHTTP(customResponseWriter, r)
	})
}
