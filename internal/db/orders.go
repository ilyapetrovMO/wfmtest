package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderModel struct {
	DB *pgxpool.Pool
}

type Order struct {
	OrderId   int64 `json:"order_id"`
	UserId    int64 `json:"user_id"`
	ProductId int64 `json:"product_id"`
	Amount    int64 `json:"amount"`
	CreatedAt int64 `json:"created_at"`
}

func (o *OrderModel) CreateOrder(ctx context.Context, user_id, product_id, amount int64, created_at time.Time) (*Order, error) {
	sql := "insert into orders(user_id, product_id, amount, created_at) values ($1, $2, $3) returning order_id, user_id, product_id, amount, created_at"
	row := o.DB.QueryRow(ctx, sql, user_id, product_id, amount, created_at)

	order := &Order{}
	err := row.Scan(&order.OrderId, &order.UserId, &order.ProductId, &order.Amount, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderModel) GetOrders(ctx context.Context) ([]*Order, error) {
	sql := "select order_id, user_id, product_id, amount, created_at from orders where deleted_at=NULL"
	rows, err := o.DB.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	orders := []*Order{}
	for rows.Next() {
		or := &Order{}
		err := rows.Scan(&or.OrderId, &or.UserId, &or.ProductId, &or.Amount, &or.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, or)
	}

	return orders, nil
}

func (o *OrderModel) GetOrdersForUser(ctx context.Context, userId int64) ([]*Order, error) {
	sql := "select order_id, user_id, product_id, amount, created_at from orders where user_id=$1 and deleted_at=NULL"
	rows, err := o.DB.Query(ctx, sql, userId)
	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	orders := []*Order{}
	for rows.Next() {
		or := &Order{}
		err := rows.Scan(&or.OrderId, &or.UserId, &or.ProductId, &or.Amount, &or.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, or)
	}

	return orders, nil
}

func (o *OrderModel) GetOrderById(ctx context.Context, orderId int64) (*Order, error) {
	sql := "select order_id, user_id, product_id, amount, created_at, deleted_at from orders where order_id=$1 and deleted_at=NULL"
	row := o.DB.QueryRow(ctx, sql, orderId)

	or := &Order{}
	err := row.Scan(&or.OrderId, &or.UserId, &or.ProductId, &or.Amount, &or.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	return or, nil
}

func (o *OrderModel) CancellOrder(ctx context.Context, orderId int64, deletedAt time.Time) error {
	sql := "update orders set deleted_at=$1 where order_id=$2"
	tag, err := o.DB.Exec(ctx, sql, deletedAt, orderId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrRecordNotFound
	}

	return nil
}
