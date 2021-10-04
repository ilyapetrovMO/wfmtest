package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	Product_id  int64  `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductModel struct {
	DB *pgxpool.Pool
}

func (p *ProductModel) CreateProduct(ctx context.Context, name, description string) (*Product, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	row := p.DB.QueryRow(ctx, "insert into products(name, description) values ($1, $2) returning product_id, name, description", name, description)
	newpr := &Product{}
	err := row.Scan(&newpr.Product_id, &newpr.Name, &newpr.Description)
	if err != nil {
		return nil, err
	}

	return newpr, nil
}
