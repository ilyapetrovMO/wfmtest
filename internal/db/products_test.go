package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestCreateProduct(t *testing.T) {
	var connstring = "postgres://postgres:postgres@127.0.0.1:5432/wfmtest"

	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	want := &Product{
		Name:        "productTest",
		Description: "descriptionTest",
	}

	pm := &ProductModel{dbpool}
	got, err := pm.CreateProduct(context.Background(), want.Name, want.Description)
	unexpectedErr(t, err)

	want.Product_id = got.Product_id
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}

	dbpool.Exec(context.Background(), "delete from products where product_id=$1", got.Product_id)
}
