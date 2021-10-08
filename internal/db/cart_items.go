package db

import (
	"context"
	"errors"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CartItem struct {
	CartItemId  int64     `json:"cart_item_id"`
	CartId      int64     `json:"cart_id"`
	Name        string    `json:"product_name"`
	Description string    `json:"product_description"`
	ProductId   int64     `json:"product_id"`
	Amount      int64     `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int64     `json:"-"`
}

type CartItemModel struct {
	DB *pgxpool.Pool
}

func ValidateCartItem(v *validator.Validator, cartItem *CartItem) {
	v.Check(cartItem.CartId > 0, "cart_id", "must be valid")
	v.Check(cartItem.ProductId > 0, "product_id", "must be valid")
	v.Check(cartItem.Amount > 0, "amount", "must be non zero positive number")
}

func (c *CartItemModel) Create(cartItem *CartItem) error {
	sql := `
	INSERT INTO cart_items (cart_id, product_id, amount, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING cart_item_id, cart_id, product_id, amount, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := c.DB.QueryRow(ctx, sql, cartItem.CartId, cartItem.ProductId, cartItem.Amount, time.Now()).Scan(
		&cartItem.CartItemId,
		&cartItem.CartId,
		&cartItem.ProductId,
		&cartItem.Amount,
		&cartItem.CreatedAt,
		&cartItem.Version,
	)
	if err != nil {
		if isPgDuplicateError(err) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (c *CartItemModel) GetById(cartItemId int64) (*CartItem, int, error) {
	sql := `
	SELECT CI.cart_item_id, CI.cart_id, CI.product_id, CI.created_at, CI.version, CI.amount, C.user_id
	FROM cart_items as CI
	JOIN carts as C
	ON C.cart_id = CI.cart_id
	WHERE CI.cart_item_id = $1 AND CI.deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var userId int
	cartItem := &CartItem{}
	err := c.DB.QueryRow(ctx, sql, cartItemId).Scan(
		&cartItem.CartItemId,
		&cartItem.CartId,
		&cartItem.ProductId,
		&cartItem.CreatedAt,
		&cartItem.Version,
		&cartItem.Amount,
		&userId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrRecordNotFound
		}
		return nil, 0, err
	}

	return cartItem, userId, nil
}

func (c *CartItemModel) GetAllByUserId(userId int64) ([]*CartItem, error) {
	sql := `
	SELECT CI.cart_item_id, CI.cart_id, CI.product_id, CI.created_at, CI.amount, CI.version, P.name, P.description
	FROM cart_items AS CI
	JOIN products AS P
	ON P.product_id = CI.product_id
	JOIN carts AS C
	ON C.cart_id = CI.cart_id
	WHERE C.user_id = $1 AND CI.deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := c.DB.Query(ctx, sql, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*CartItem{}, nil
		}
		return nil, err
	}

	items := []*CartItem{}
	for rows.Next() {
		item := &CartItem{}
		err := rows.Scan(
			&item.CartItemId,
			&item.CartId,
			&item.ProductId,
			&item.CreatedAt,
			&item.Amount,
			&item.Version,
			&item.Name,
			&item.Description,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (c *CartItemModel) Delete(item *CartItem) error {
	sql := `
	UPDATE cart_items
	SET deleted_at = $1, version = version + 1
	WHERE cart_item_id = $2 AND version = $3 AND deleted_at IS NULL
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := c.DB.QueryRow(ctx, sql, time.Now(), item.CartItemId, item.Version).Scan(&item.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}
