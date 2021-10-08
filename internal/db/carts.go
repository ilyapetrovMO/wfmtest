package db

import (
	"context"
	"errors"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Cart struct {
	CartId    int64     `json:"cart_id"`
	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"-"`
}

type CartModel struct {
	DB *pgxpool.Pool
}

func ValidateCart(v *validator.Validator, cart *Cart) {
	v.Check(cart.CartId > 0, "cart_id", "must be valid")
	v.Check(cart.UserId > 0, "user_id", "must be valid")
}

func (c *CartModel) Create(cart *Cart) error {
	sql := `
	INSERT INTO carts (user_id, created_at)
	VALUES ($1, $2)
	RETURNING cart_id, user_id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := c.DB.QueryRow(ctx, sql, cart.UserId, time.Now()).Scan(
		&cart.CartId,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.Version,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartModel) GetByUserId(userId int64) (*Cart, error) {
	sql := `
	SELECT cart_id, user_id, created_at, version
	FROM carts
	WHERE user_id = $1 AND deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cart := &Cart{}
	err := c.DB.QueryRow(ctx, sql, userId).Scan(
		&cart.CartId,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return cart, nil
}

func (c *CartModel) GetById(cartId int64) (*Cart, error) {
	sql := `
	SELECT cart_id, user_id, created_at, version
	FROM carts
	WHERE cart_id = $1 AND deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cart := &Cart{}
	err := c.DB.QueryRow(ctx, sql, cartId).Scan(
		&cart.CartId,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return cart, nil
}

func (c *CartModel) GetAll() ([]*Cart, error) {
	sql := `
	SELECT cart_id, user_id, created_at, version
	FROM carts
	WHERE deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := c.DB.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	carts := []*Cart{}
	for rows.Next() {
		cart := &Cart{}
		err = rows.Scan(
			&cart.CartId,
			&cart.UserId,
			&cart.CreatedAt,
			&cart.Version,
		)
		if err != nil {
			return nil, err
		}

		carts = append(carts, cart)
	}

	return carts, nil
}

func (c *CartModel) Delete(cart *Cart) error {
	sql := `
	UPDATE carts
	SET deleted_at = $1, version = version+1
	WHERE cart_id = $2 AND version = $3 AND deleted_at IS NULL
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := c.DB.QueryRow(ctx, sql, time.Now(), cart.CartId, cart.Version).Scan(&cart.Version)
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
