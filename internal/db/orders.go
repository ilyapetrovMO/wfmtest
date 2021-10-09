package db

import (
	"context"
	"errors"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderModel struct {
	DB *pgxpool.Pool
}

type Order struct {
	OrderId   int64     `json:"order_id"`
	UserId    int64     `json:"user_id"`
	ProductId int64     `json:"product_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"-"`
}

func ValidateOrder(v *validator.Validator, order *Order) {
	v.Check(order.UserId > 0, "user_id", "must be valid")
	v.Check(order.ProductId > 0, "product_id", "must be valid")
	v.Check(order.Amount > 0, "amount", "must be non zero positive number")
}

func (o *OrderModel) CreateTx(cartItem *CartItem, product *Product, userId int64) (*Order, error) {
	sqlProduct := `
	UPDATE products
	SET in_storage = in_storage-$1, version = version+1
	WHERE product_id = $2 AND deleted_at IS NULL AND version = $3`

	sqlOrder := `
	INSERT INTO orders(user_id, product_id, amount, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING order_id, user_id, product_id, amount, created_at, version`

	sqlCartItem := `
	UPDATE cart_items
	SET deleted_at = $1, version = version + 1
	WHERE cart_item_id = $2 AND version = $3 AND deleted_at IS NULL
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := o.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, sqlProduct, cartItem.Amount, product.ProductId, product.Version)
	if err != nil {
		return nil, err
	}

	order := &Order{}
	err = tx.QueryRow(ctx, sqlOrder, userId, product.ProductId, cartItem.Amount, time.Now()).Scan(
		&order.OrderId,
		&order.UserId,
		&order.ProductId,
		&order.Amount,
		&order.CreatedAt,
		&order.Version,
	)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, sqlCartItem, time.Now(), cartItem.CartItemId, cartItem.Version).Scan(&cartItem.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	tx.Commit(ctx)
	return order, nil
}

func (o *OrderModel) Create(order *Order) error {
	sql := `
	INSERT INTO orders(user_id, product_id, amount, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING order_id, user_id, product_id, amount, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.DB.QueryRow(ctx, sql, order.UserId, order.ProductId, order.Amount, time.Now()).Scan(
		&order.OrderId,
		&order.UserId,
		&order.ProductId,
		&order.Amount,
		&order.CreatedAt,
		&order.Version,
	)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderModel) GetAll() ([]*Order, error) {
	sql := `SELECT
	order_id, user_id, product_id, amount, created_at, version
	FROM orders
	WHERE deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := o.DB.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		return []*Order{}, nil
	}
	if err != nil {
		return nil, err
	}

	orders := []*Order{}
	for rows.Next() {
		or := &Order{}
		err := rows.Scan(
			&or.OrderId,
			&or.UserId,
			&or.ProductId,
			&or.Amount,
			&or.CreatedAt,
			&or.Version,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, or)
	}

	return orders, nil
}

func (o *OrderModel) GetWithUserId(userId int64) ([]*Order, error) {
	sql := `
	SELECT order_id, user_id, product_id, amount, created_at, version
	FROM orders
	WHERE user_id=$1 AND deleted_at IS NULL`

	rows, err := o.DB.Query(context.Background(), sql, userId)
	if err == pgx.ErrNoRows {
		return []*Order{}, nil
	}
	if err != nil {
		return nil, err
	}

	orders := []*Order{}
	for rows.Next() {
		or := &Order{}
		err := rows.Scan(
			&or.OrderId,
			&or.UserId,
			&or.ProductId,
			&or.Amount,
			&or.CreatedAt,
			&or.Version,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, or)
	}

	return orders, nil
}

func (o *OrderModel) GetById(orderId int64) (*Order, error) {
	sql := `
	SELECT order_id, user_id, product_id, amount, created_at, version
	FROM orders
	WHERE order_id=$1 and deleted_at IS NULL`

	or := &Order{}
	err := o.DB.QueryRow(context.Background(), sql, orderId).Scan(
		&or.OrderId,
		&or.UserId,
		&or.ProductId,
		&or.Amount,
		&or.CreatedAt,
		&or.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return or, nil
}

func (o *OrderModel) Delete(order *Order) error {
	sql := `
	UPDATE orders
	SET deleted_at = $1, version = version+1
	WHERE order_id = $2 AND version = $3 AND deleted_at IS NULL
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := o.DB.QueryRow(ctx, sql, time.Now(), order.OrderId, order.Version).Scan(&order.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}
