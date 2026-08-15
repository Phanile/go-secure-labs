package domain

import "time"

type Order struct {
	Id         int64
	UserId     int64
	AddressId  int64
	Status     string
	TotalPrice int64 // better to use custom Money type
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
