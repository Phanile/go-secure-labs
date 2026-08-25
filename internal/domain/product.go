package domain

import (
	"fmt"
	"time"

	core_errors "github.com/Phanile/go-secure-labs/internal/errors"
)

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

func NewProductUninitialized(name, description string,
	price int64,
	stockQuantity, categoryId int,
	imageUrl string,
) Product {
	return Product{
		Id:            UninitializedID,
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

func NewProduct(id int64,
	name, description string,
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

func (p Product) Validate() error {
	if p.Price <= 0 {
		return fmt.Errorf("validate product price %d must be > 0: %w", p.Price, core_errors.ErrInvalidArgument)
	}

	if p.StockQuantity < 0 {
		return fmt.Errorf("validate product StockQuantity %d must be >= 0: %w", p.StockQuantity, core_errors.ErrInvalidArgument)
	}

	return nil
}
