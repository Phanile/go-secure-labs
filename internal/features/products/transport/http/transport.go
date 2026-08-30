package products_transport_http

import (
	"context"
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/domain"
	products_service "github.com/Phanile/go-secure-labs/internal/features/products/service"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
)

type ProductsHTTPHandler struct {
	productsService ProductsService
}

type ProductsService interface {
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	GetProducts(ctx context.Context, filter products_service.ProductFilter) ([]domain.Product, error)
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
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: h.GetProducts,
		},
	}
}

func (h *ProductsHTTPHandler) VulnerableRoutes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/vulnerable/products",
			Handler: h.GetProducts,
		},
	}
}
