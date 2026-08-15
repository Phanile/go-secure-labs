package domain

import "time"

type OrderItem struct {
	Id        int64
	OrderId   int64
	ProductId int64
	Quantity  int
	UnitPrice int64 // better to use custom Money type
	CreatedAt time.Time
}

func NewOrderItem(id, orderId, productId int64, quantity int, unitPrice int64) OrderItem {
	return OrderItem{
		Id:        id,
		OrderId:   orderId,
		ProductId: productId,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		CreatedAt: time.Now(),
	}
}
