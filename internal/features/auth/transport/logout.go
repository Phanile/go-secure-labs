package auth_transport_http

import (
	"net/http"
	"time"

	internal_errors "github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_middleware "github.com/Phanile/go-secure-labs/internal/transport/http/middleware"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
)

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userId, ok := http_middleware.UserIdFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(internal_errors.ErrUnauthorized, "user id not found in context")
		return
	}

	err := h.authService.Logout(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
}
