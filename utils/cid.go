package utils

import (
	"context"
	"fmt"
	"golang-boiler-plate/constants"
	"math/rand"
	"net/http"
	"time"
)

func GetCorrelationIDFromRequest(r *http.Request) string {
	cid := r.Header.Get(string(constants.HEADER_CORRELATION_ID))
	if len(cid) <= 0 {
		return GenerateCorrelationID()
	}

	return cid
}

func GenerateCorrelationID() string {
	rangeLower := 100000
	rangeUpper := 999999
	rand.NewSource(time.Now().UnixNano())
	CID := fmt.Sprintf("CID-%s-%d", getOneLineCurrentTimeString(), rangeLower+rand.Intn(rangeUpper-rangeLower+1))
	return CID
}

func GetNewContextWithCID(r *http.Request, CIDstr string) (newCtx context.Context) {
	ctx := r.Context()
	return context.WithValue(ctx, constants.CORRELATION_ID_CONTEXT_KEY, CIDstr)
}

func GetCorrelationIDFromContext(ctx context.Context) string {
	if ctx == nil || ctx.Value(constants.CORRELATION_ID_CONTEXT_KEY) == nil {
		return GenerateCorrelationID()
	}

	return ctx.Value(constants.CORRELATION_ID_CONTEXT_KEY).(string)
}

func getOneLineCurrentTimeString() string {
	t := time.Now()
	return fmt.Sprintf("%v%.2v%.2v%.2v%.2v%.2v", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
}
