package products_transport_http

import (
	"encoding/json"
	"net/http"
)

type CreateProductRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"` // better to use custom Money type
	StockQuantity int    `json:"stock_quantity"`
	CategoryId    int    `json:"category_id"`
	ImageUrl      string `json:"image_url"`
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
	var request CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

	}
}
