package response

import (
	"context"
	"encoding/json"
	"golang-boiler-plate/utils/logger"
	"net/http"
)

// Response is response struct
type Response struct {
	w       http.ResponseWriter
	ctx     context.Context
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}

func New(ctx context.Context, w http.ResponseWriter) *Response {
	return &Response{
		ctx:    ctx,
		w:      w,
		Status: http.StatusOK,
	}
}

func (res *Response) SetCode(code int) *Response {
	res.Status = code
	return res
}

func (res *Response) SetData(data interface{}) *Response {
	res.Data = data
	return res
}

func (res *Response) SetMessage(message string) *Response {
	res.Message = message
	return res
}

func (res *Response) Write() {
	w := res.w
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(res.Status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Ctx(res.ctx).StackTrace(err).Error("Error encoding response")
	}
}
