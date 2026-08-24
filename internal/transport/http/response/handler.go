package http_response

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.w.WriteHeader(statusCode)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("PanicResponse encoder err: ", zap.Error(err))
	}
}
