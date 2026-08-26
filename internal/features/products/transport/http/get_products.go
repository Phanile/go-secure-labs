package products_transport_http

import (
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/domain"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
)

type GetProductsResponse struct {
	products []ProductDTO
}

type ProductDTO struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"` // better to use custom Money type
	StockQuantity int    `json:"stock_quantity"`
	CategoryId    int    `json:"category_id"`
	ImageUrl      string `json:"image_url"`
}

func (h *ProductsHTTPHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	products, err := h.productsService.GetProducts(ctx)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get products")
		return
	}

	productDTOs := make([]ProductDTO, 0)

	for _, p := range products {
		productDTOs = append(productDTOs, domainToDTO(p))
	}

	responseHandler.JSONResponse(productDTOs, http.StatusOK)
}

func domainToDTO(p domain.Product) ProductDTO {
	return ProductDTO{
		Id:            p.Id,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		CategoryId:    p.CategoryId,
		ImageUrl:      p.ImageUrl,
	}
}
