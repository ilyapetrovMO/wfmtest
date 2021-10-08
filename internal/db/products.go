package db

import (
	"context"
	"errors"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	ProductId   int64     `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	InStorage   int64     `json:"in_storage"`
	Version     int64     `json:"-"`
}

type ProductModel struct {
	DB *pgxpool.Pool
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(product.InStorage >= 0, "int_storage", "must be a positive number")
}

func (p *ProductModel) Create(product *Product) error {
	sql := `
	INSERT INTO products(name, description, created_at, in_storage)
	VALUES ($1, $2, $3, $4)
	RETURNING product_id, name, description, created_at, in_storage, version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	args := []interface{}{
		product.Name,
		product.Description,
		time.Now(),
		product.InStorage,
	}

	row := p.DB.QueryRow(ctx, sql, args...)
	err := row.Scan(
		&product.ProductId,
		&product.Name,
		&product.Description,
		&product.CreatedAt,
		&product.InStorage,
		&product.Version,
	)
	if err != nil {
		if isPgDuplicateError(err) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (p *ProductModel) GetAll() ([]*Product, error) {
	sql := `
	SELECT product_id, name, description, created_at, in_storage, version
	FROM products
	WHERE deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := p.DB.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		return []*Product{}, nil
	}
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		pr := &Product{}
		err = rows.Scan(
			&pr.ProductId,
			&pr.Name,
			&pr.Description,
			&pr.CreatedAt,
			&pr.InStorage,
			&pr.Version,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, pr)
	}

	return products, nil
}

func (p *ProductModel) GetById(productId int64) (*Product, error) {
	sql := `
	SELECT product_id, name, description, created_at, in_storage, version
	FROM products
	WHERE product_id = $1 AND deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	product := &Product{}

	err := p.DB.QueryRow(ctx, sql, productId).Scan(
		&product.ProductId,
		&product.Name,
		&product.Description,
		&product.CreatedAt,
		&product.InStorage,
		&product.Version,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductModel) Update(product *Product) error {
	sql := `
	UPDATE products
	SET name = $1, description = $2, in_storage = $3, version = version + 1
	WHERE product_id = $4 AND version = $5 AND deleted_at IS NULL
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	args := []interface{}{
		product.Name,
		product.Description,
		product.InStorage,
		product.ProductId,
		product.Version,
	}

	err := p.DB.QueryRow(ctx, sql, args...).Scan(&product.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (p *ProductModel) Delete(product *Product) error {
	sql := `
	UPDATE products
	SET deleted_at = $1, version=version+1
	WHERE product_id = $2 AND version = $3 AND deleted_at IS NULL
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := p.DB.QueryRow(ctx, sql, time.Now(), product.ProductId, product.Version).Scan(&product.Version)
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
