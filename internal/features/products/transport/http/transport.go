package products_transport_http

import (
	"net/http"

	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
)

type ProductsHTTPHandler struct {
	productsService ProductsService
}

type ProductsService interface {
}

func NewProductsHTTPHandler(productsService ProductsService) *ProductsHTTPHandler {
	return &ProductsHTTPHandler{
		productsService: productsService,
	}
}

func (h *ProductsHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/products",
			Handler: h.CreateProduct,
		},
	}
}
