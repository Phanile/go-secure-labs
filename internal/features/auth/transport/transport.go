package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	Register(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, domain.Session, error)
	Login(ctx context.Context, loginUserDTO dto.LoginUserDTO) (domain.User, domain.Session, error)
	Me(ctx context.Context, userId int64) (domain.User, error)
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodGet,
			Path:    "/auth/me",
			Handler: h.Me,
		},
	}
}

// func (h *UserHTTPHandler) VulnerableRoutes() []http_server.Route {
// 	return []http_server.Route{
// 		{
// 			Method:  http.MethodGet,
// 			Path:    "/vulnerable/users",
// 			Handler: h.GetUsers,
// 		},
// 	}
// }
