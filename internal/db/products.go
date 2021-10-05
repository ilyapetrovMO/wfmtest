package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	ProductId   int64     `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductModel struct {
	DB *pgxpool.Pool
}

func (p *ProductModel) CreateProduct(ctx context.Context, name, description string, date time.Time) (*Product, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	row := p.DB.QueryRow(ctx, "insert into products(name, description, created_at) values ($1, $2, $3) returning product_id, name, description, created_at", name, description, date)
	newpr := &Product{}
	err := row.Scan(&newpr.ProductId, &newpr.Name, &newpr.Description, &newpr.CreatedAt)
	if err != nil {
		return nil, err
	}

	return newpr, nil
}

func (p *ProductModel) GetProducts(ctx context.Context) ([]*Product, error) {
	rows, err := p.DB.Query(ctx, "select product_id, name, description from products")
	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		pr := &Product{}
		err = rows.Scan(&pr.ProductId, &pr.Name, &pr.Description)
		if err != nil {
			return nil, err
		}

		products = append(products, pr)
	}

	return products, nil
}
