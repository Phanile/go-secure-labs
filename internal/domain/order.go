package domain

import "time"

type Order struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	AddressId  int64     `json:"address_id"`
	Status     string    `json:"status"`
	TotalPrice int64     `json:"total_price"` // better to use custom Money type
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewOrder(id, userId, addressId int64, status string, totalPrice int64) Order {
	return Order{
		Id:         id,
		UserId:     userId,
		AddressId:  addressId,
		Status:     status,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
