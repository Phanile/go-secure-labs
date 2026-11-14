package orders_service

import (
	"context"
	"github.com/Phanile/go-secure-labs/internal/domain"
)

type CreateOrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
type CreateOrderInput struct {
	AddressID int64             `json:"address_id,omitempty"`
	Address   AddressInput      `json:"address"`
	Items     []CreateOrderItem `json:"items"`
}
type AddressInput struct {
	Country    string `json:"country"`
	City       string `json:"city"`
	Street     string `json:"street"`
	House      string `json:"house"`
	Apartment  string `json:"apartment,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}
type Order struct {
	domain.Order
	Items []OrderItem `json:"items"`
}
type OrderItem struct {
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}
type Repository interface {
	CreateOrder(context.Context, int64, CreateOrderInput) (Order, error)
	GetOrders(context.Context, int64, bool) ([]Order, error)
	GetOrder(context.Context, int64, int64, bool) (Order, error)
}
type Service struct{ repo Repository }

func New(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, userID int64, in CreateOrderInput) (Order, error) {
	return s.repo.CreateOrder(ctx, userID, in)
}
func (s *Service) List(ctx context.Context, userID int64, vulnerable bool) ([]Order, error) {
	return s.repo.GetOrders(ctx, userID, vulnerable)
}
func (s *Service) Get(ctx context.Context, userID, orderID int64, vulnerable bool) (Order, error) {
	return s.repo.GetOrder(ctx, userID, orderID, vulnerable)
}
