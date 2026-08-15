package domain

import "time"

type Address struct {
	Id         int64
	UserId     int64
	Country    string
	City       string
	Street     string
	House      string
	Apartment  *string
	PostalCode *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewAddress(id, userId int64, country, city, street, house string, apartment, postalCode *string) Address {
	return Address{
		Id:         id,
		UserId:     userId,
		Country:    country,
		City:       city,
		Street:     street,
		House:      house,
		Apartment:  apartment,
		PostalCode: postalCode,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
