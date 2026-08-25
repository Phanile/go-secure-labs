package products_postgres_repository

type ProductModel struct {
	Id            int64
	Name          string
	Description   string
	Price         int64 // better to use custom Money type
	StockQuantity int
	CategoryId    int
	ImageUrl      string
}
