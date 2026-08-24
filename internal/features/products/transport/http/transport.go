package products_transport_http

import (
	"context"
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/domain"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
)

type ProductsHTTPHandler struct {
	productsService ProductsService
}

type ProductsService interface {
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
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
