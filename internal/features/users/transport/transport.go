package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/dto"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(ctx context.Context, createUserDTO dto.CreateUserDTO) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: usersService,
	}
}

func (h *UserHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: h.CreateUser,
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
