package domain

import "time"

type Product struct {
	Id            int64
	Name          string
	Description   string
	Price         int64 // better to use custom Money type
	StockQuantity int
	CategoryId    int
	ImageUrl      string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewProduct(id int64, name, description string,
	price int64,
	stockQuantity, categoryId int,
	imageUrl string,
) Product {
	return Product{
		Id:            id,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		CategoryId:    categoryId,
		ImageUrl:      imageUrl,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
