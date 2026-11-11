package http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/Phanile/go-secure-labs/internal/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(l *logger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: l,
		w:   w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("json response encoder err: ", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(e error, msg string) {
	var (
		statusCode int
		logFunc    func(msg string, fields ...zap.Field)
	)

	switch {
	case errors.Is(e, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(e, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(e, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(e, core_errors.ErrPasswordMismatch):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(e, core_errors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn

	case errors.Is(e, core_errors.ErrTokenExpired):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(e))

	h.errorResponse(msg, statusCode, e)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.errorResponse(msg, statusCode, err)
}

func (h *HTTPResponseHandler) errorResponse(
	msg string,
	statusCode int,
	err error,
) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
