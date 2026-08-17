package products_transport_http

import "net/http"

type CreateProductRequest struct {
	Name          string
	Description   string
	Price         int64 // better to use custom Money type
	StockQuantity int
	CategoryId    int
	ImageUrl      string
}

func (h *ProductsHTTPHandler) CreateProduct(rw http.ResponseWriter, r *http.Request) {

}
