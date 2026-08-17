package products_transport_http

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
