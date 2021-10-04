package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderModel struct {
	DB *pgxpool.Pool
}

type Order struct {
	Order_id   int64 `json:"order_id"`
	User_id    int64 `json:"user_id"`
	Product_id int64 `json:"product_id"`
	Amount     int64 `json:"amount"`
}

func (o *OrderModel) CreateOrder(ctx context.Context, user_id, product_id, amount int64) (*Order, error) {
	sql := "insert into orders(user_id, product_id, amount) values ($1, $2, $3) returning order_id, user_id, product_id, amount"
	row := o.DB.QueryRow(ctx, sql, user_id, product_id, amount)

	order := &Order{}
	err := row.Scan(&order.Order_id, &order.User_id, &order.Product_id, &order.Amount)
	if err != nil {
		return nil, err
	}

	return order, nil
}
