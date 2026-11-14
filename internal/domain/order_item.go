package domain

import "time"

type OrderItem struct {
	Id        int64     `json:"id"`
	OrderId   int64     `json:"order_id"`
	ProductId int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"` // better to use custom Money type
	CreatedAt time.Time `json:"created_at"`
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
