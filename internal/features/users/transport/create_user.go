package users_transport_http

import (
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/dto"
	"github.com/Phanile/go-secure-labs/internal/errors"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_request "github.com/Phanile/go-secure-labs/internal/transport/http/request"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
)

type CreateUserRequest struct {
	Email          string `json:"email"           validate:"required,email"`
	Password       string `json:"password"        validate:"required"`
	RepeatPassword string `json:"repeat_password" validate:"required"`
}

type CreateUserResponse struct {
	Id int64 `json:"id"`
}

func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request CreateUserRequest

	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	if request.Password != request.RepeatPassword {
		responseHandler.ErrorResponse(errors.ErrPasswordMismatch, "password and repeat password do not match")
		return
	}

	createUserDTO := dto.CreateUserDTO{
		Email:    request.Email,
		Password: request.Password,
	}

	user, err := h.userService.CreateUser(ctx, createUserDTO)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create new user")
		return
	}

	resp := CreateUserResponse{
		Id: user.Id,
	}

	responseHandler.JSONResponse(resp, http.StatusCreated)
}
