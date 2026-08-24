package http_server

import (
	"fmt"
	"net/http"
)

type APIVersion string

const (
	APIVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
}

func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (router *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		router.Handle(pattern, route.Handler)
	}
}
