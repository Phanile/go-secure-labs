package domain

import "time"

type Basket struct {
	Id        int64
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBasket(id, userId int64) Basket {
	return Basket{
		Id:        id,
		UserId:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
