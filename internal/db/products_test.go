package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var connstring = "postgres://postgres:postgres@127.0.0.1:5432/wfmtest"

func TestCreateProduct(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	want := &Product{
		Name:        "productTest",
		Description: "descriptionTest",
		CreatedAt:   time.Now(),
	}

	pm := &ProductModel{dbpool}
	got, err := pm.CreateProduct(context.Background(), want.Name, want.Description, want.CreatedAt)
	unexpectedErr(t, err)

	dbpool.Exec(context.Background(), "delete from products where product_id=$1", got.ProductId)

	if got.Name != want.Name {
		t.Errorf("name: got %s want %s", got.Name, want.Name)
	}
	if got.Description != want.Description {
		t.Errorf("name: got %s want %s", got.Description, want.Description)
	}
	if got.CreatedAt.Day() != want.CreatedAt.Day() {
		t.Errorf("name: got %s want %s", got.CreatedAt, want.CreatedAt)
	}

	dbpool.Exec(context.Background(), "delete from products where product_id=$1", got.ProductId)
}

func TestGetProducts(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	want := 3

	pm := &ProductModel{dbpool}
	got, err := pm.GetProducts(context.Background())
	unexpectedErr(t, err)

	if len(got) != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}
