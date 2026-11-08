package products_transport_http

import (
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_request "github.com/Phanile/go-secure-labs/internal/transport/http/request"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
)

type CreateProductRequest struct {
	Name          string `json:"name"           validate:"required"`
	Description   string `json:"description"    validate:"required"`
	Price         int64  `json:"price"          validate:"required"` // better to use custom Money type
	StockQuantity int    `json:"stock_quantity" validate:"required"`
	CategoryId    int    `json:"category_id"    validate:"required"`
	ImageUrl      string `json:"image_url"      validate:"required"`
}

type CreateProductResponse struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"` // better to use custom Money type
	StockQuantity int    `json:"stock_quantity"`
	CategoryId    int    `json:"category_id"`
	ImageUrl      string `json:"image_url"`
}

func (h *ProductsHTTPHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request CreateProductRequest

	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	product, err := h.productsService.CreateProduct(ctx, domainFromDTO(request))

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create new product")
		return
	}

	resp := CreateProductResponse{
		Id:            product.Id,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryId:    product.CategoryId,
		ImageUrl:      product.ImageUrl,
	}

	responseHandler.JSONResponse(resp, http.StatusCreated)
}

func domainFromDTO(dto CreateProductRequest) domain.Product {
	return domain.NewProductUninitialized(
		dto.Name,
		dto.Description,
		dto.Price,
		dto.StockQuantity,
		dto.CategoryId,
		dto.ImageUrl,
	)
}
