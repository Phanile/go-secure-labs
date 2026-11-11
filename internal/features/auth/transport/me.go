package auth_transport_http

import (
	"net/http"

	internal_errors "github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_middleware "github.com/Phanile/go-secure-labs/internal/transport/http/middleware"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
)

type MeResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *AuthHTTPHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userId, ok := http_middleware.UserIdFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(internal_errors.ErrUnauthorized, "user id not found in context")
		return
	}

	user, err := h.authService.Me(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(internal_errors.ErrUnauthorized, "user id not found in context")
		return
	}

	resp := MeResponse{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	responseHandler.JSONResponse(resp, http.StatusOK)
}
