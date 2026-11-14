package orders_postgres_repository

import (
	"context"
	"fmt"
	core_errors "github.com/Phanile/go-secure-labs/internal/errors"
	orders_service "github.com/Phanile/go-secure-labs/internal/features/orders/service"
	core "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
)

type Repository struct{ pool core.Pool }

func New(pool core.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) GetOrders(ctx context.Context, userID int64, vulnerable bool) ([]orders_service.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `SELECT id,user_id,address_id,status,total_price,created_at,updated_at FROM appleshop.orders WHERE user_id=$1 ORDER BY created_at DESC`
	if vulnerable {
		query = `SELECT id,user_id,address_id,status,total_price,created_at,updated_at FROM appleshop.orders ORDER BY created_at DESC`
	}
	var rows pgx.Rows
	var err error
	if vulnerable {
		rows, err = r.pool.Query(ctx, query)
	} else {
		rows, err = r.pool.Query(ctx, query, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()
	result := make([]orders_service.Order, 0)
	for rows.Next() {
		var o orders_service.Order
		if err := rows.Scan(&o.Id, &o.UserId, &o.AddressId, &o.Status, &o.TotalPrice, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, rows.Err()
}
func (r *Repository) GetOrder(ctx context.Context, userID, orderID int64, vulnerable bool) (orders_service.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `SELECT id,user_id,address_id,status,total_price,created_at,updated_at FROM appleshop.orders WHERE id=$1 AND user_id=$2`
	args := []any{orderID, userID}
	if vulnerable {
		query = `SELECT id,user_id,address_id,status,total_price,created_at,updated_at FROM appleshop.orders WHERE id=$1`
		args = []any{orderID}
	}
	var o orders_service.Order
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&o.Id, &o.UserId, &o.AddressId, &o.Status, &o.TotalPrice, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return o, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
		}
		return o, fmt.Errorf("get order: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT oi.product_id,p.name,oi.quantity,oi.unit_price FROM appleshop.order_items oi JOIN appleshop.products p ON p.id=oi.product_id WHERE oi.order_id=$1 ORDER BY oi.id`, orderID)
	if err != nil {
		return o, err
	}
	defer rows.Close()
	o.Items = make([]orders_service.OrderItem, 0)
	for rows.Next() {
		var i orders_service.OrderItem
		if err := rows.Scan(&i.ProductID, &i.Name, &i.Quantity, &i.UnitPrice); err != nil {
			return o, err
		}
		o.Items = append(o.Items, i)
	}
	return o, rows.Err()
}
func (r *Repository) CreateOrder(ctx context.Context, userID int64, in orders_service.CreateOrderInput) (orders_service.Order, error) {
	if (in.AddressID <= 0 && (in.Address.Country == "" || in.Address.City == "" || in.Address.Street == "" || in.Address.House == "")) || len(in.Items) == 0 {
		return orders_service.Order{}, fmt.Errorf("address and items are required: %w", core_errors.ErrInvalidArgument)
	}
	for _, item := range in.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return orders_service.Order{}, fmt.Errorf("invalid order item: %w", core_errors.ErrInvalidArgument)
		}
	}
	// Merge duplicate product IDs before locking and checking inventory.
	quantities := make(map[int64]int, len(in.Items))
	for _, item := range in.Items {
		quantities[item.ProductID] += item.Quantity
	}
	in.Items = in.Items[:0]
	for id, quantity := range quantities {
		in.Items = append(in.Items, orders_service.CreateOrderItem{ProductID: id, Quantity: quantity})
	}
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return orders_service.Order{}, err
	}
	defer tx.Rollback(ctx)
	if in.AddressID > 0 {
		var ownsAddress bool
		if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM appleshop.addresses WHERE id=$1 AND user_id=$2)`, in.AddressID, userID).Scan(&ownsAddress); err != nil {
			return orders_service.Order{}, err
		}
		if !ownsAddress {
			return orders_service.Order{}, fmt.Errorf("address not found: %w", core_errors.ErrNotFound)
		}
	} else {
		err = tx.QueryRow(ctx, `INSERT INTO appleshop.addresses(user_id,country,city,street,house,apartment,postal_code) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, userID, in.Address.Country, in.Address.City, in.Address.Street, in.Address.House, in.Address.Apartment, in.Address.PostalCode).Scan(&in.AddressID)
		if err != nil {
			return orders_service.Order{}, fmt.Errorf("create shipping address: %w", err)
		}
	}
	var total int64
	for _, i := range in.Items {
		var price int64
		var stock int
		if err = tx.QueryRow(ctx, `SELECT price,stock_quantity FROM appleshop.products WHERE id=$1 AND is_active=true FOR UPDATE`, i.ProductID).Scan(&price, &stock); err != nil {
			if err == pgx.ErrNoRows {
				return orders_service.Order{}, fmt.Errorf("product unavailable: %w", core_errors.ErrNotFound)
			}
			return orders_service.Order{}, fmt.Errorf("product unavailable: %w", err)
		}
		if stock < i.Quantity {
			return orders_service.Order{}, fmt.Errorf("insufficient stock: %w", core_errors.ErrInvalidArgument)
		}
		total += price * int64(i.Quantity)
	}
	var o orders_service.Order
	if err = tx.QueryRow(ctx, `INSERT INTO appleshop.orders(user_id,address_id,status,total_price) VALUES($1,$2,'pending',$3) RETURNING id,user_id,address_id,status,total_price,created_at,updated_at`, userID, in.AddressID, total).Scan(&o.Id, &o.UserId, &o.AddressId, &o.Status, &o.TotalPrice, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return o, err
	}
	for _, i := range in.Items {
		var price int64
		if err = tx.QueryRow(ctx, `SELECT price FROM appleshop.products WHERE id=$1`, i.ProductID).Scan(&price); err != nil {
			return o, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO appleshop.order_items(order_id,product_id,quantity,unit_price) VALUES($1,$2,$3,$4)`, o.Id, i.ProductID, i.Quantity, price); err != nil {
			return o, err
		}
		if _, err = tx.Exec(ctx, `UPDATE appleshop.products SET stock_quantity=stock_quantity-$2 WHERE id=$1`, i.ProductID, i.Quantity); err != nil {
			return o, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return o, err
	}
	return r.GetOrder(ctx, userID, o.Id, false)
}
