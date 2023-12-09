package middleware

import (
	"golang-boiler-plate/utils"
	"golang-boiler-plate/utils/logger"
	"net/http"
)

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := utils.GetCorrelationIDFromRequest(r)
		ctx := utils.GetNewContextWithCID(r, cid)
		logCtx := logger.GetCorrelationIDLoggerCtx(ctx, cid)

		next.ServeHTTP(w, r.WithContext(logCtx))
	})
}
