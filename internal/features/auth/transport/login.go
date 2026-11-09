package auth_transport_http

import (
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/dto"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_request "github.com/Phanile/go-secure-labs/internal/transport/http/request"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email"           validate:"required,email"`
	Password string `json:"password"        validate:"required"`
}

type LoginResponse struct {
	Id int64 `json:"id"`
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request LoginRequest

	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	loginUserDTO := dto.LoginUserDTO{
		Email:    request.Email,
		Password: request.Password,
	}

	user, session, err := h.authService.Login(ctx, loginUserDTO)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	resp := LoginResponse{
		Id: user.Id,
	}

	responseHandler.JSONResponse(resp, http.StatusOK)
}
