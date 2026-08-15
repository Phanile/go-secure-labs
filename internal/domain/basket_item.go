package domain

import "time"

type BasketItem struct {
	Id        int64
	BasketId  int64
	ProductId int64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBasketItem(id, basketId, productId int64, quantity int) BasketItem {
	return BasketItem{
		Id:        id,
		BasketId:  basketId,
		ProductId: productId,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
