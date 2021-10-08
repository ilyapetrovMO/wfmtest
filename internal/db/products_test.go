package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4/pgxpool"
)

var connstring = "postgres://postgres:postgres@127.0.0.1:5432/wfmtest"

func TestProductsValidateProduct(t *testing.T) {
	testCases := []struct {
		desc    string
		product *Product
		want    bool
	}{
		{
			desc:    "empty name",
			product: &Product{ProductId: 1},
			want:    false,
		}, {
			desc:    "valid product",
			product: &Product{ProductId: 1, Name: "name"},
			want:    true,
		}, {
			desc:    "negative int_storage",
			product: &Product{ProductId: 1, Name: "name", InStorage: -1},
			want:    false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			v := validator.New()

			ValidateProduct(v, tC.product)

			if v.Valid() != tC.want {
				t.Errorf("got %v want %v", v.Valid(), tC.want)
			}
		})
	}
}

func TestProductsCreate(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()
	defer dbpool.Exec(context.Background(), `DELETE FROM products WHERE name LIKE '%TEST'`)

	t.Run("valid product creation", func(t *testing.T) {
		want := &Product{
			Name:        "testCreateTEST",
			Description: "valid product creation",
			CreatedAt:   time.Now(),
		}

		got := &Product{
			Name:        "testCreateTEST",
			Description: "valid product creation",
			CreatedAt:   time.Now(),
		}

		pm := &ProductModel{dbpool}
		err := pm.Create(got)
		unexpectedErr(t, err)

		if got.Name != want.Name {
			t.Errorf("name: got %s want %s", got.Name, want.Name)
		}
		if got.Description != want.Description {
			t.Errorf("name: got %s want %s", got.Description, want.Description)
		}
		if got.CreatedAt.Day() != want.CreatedAt.Day() {
			t.Errorf("name: got %s want %s", got.CreatedAt, want.CreatedAt)
		}
	})

	t.Run("try to add twice", func(t *testing.T) {
		want := &Product{
			Name:        "addTwiceTEST",
			Description: "valid product creation",
			CreatedAt:   time.Now(),
		}

		pm := &ProductModel{dbpool}
		err := pm.Create(want)
		unexpectedErr(t, err)

		err = pm.Create(want)

		if err != nil {
			switch {
			case errors.Is(err, ErrAlreadyExists):
				return
			default:
				unexpectedErr(t, err)
			}
		}

		t.Fatal("expected error")
	})
}

func TestProductsGetAll(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	want := 3

	pm := &ProductModel{dbpool}
	got, err := pm.GetAll()
	unexpectedErr(t, err)

	if len(got) != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestProductsDelete(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()
	defer dbpool.Exec(context.Background(), `DELETE FROM products WHERE name LIKE '%TEST'`)

	t.Run("try deleting nonexistant product", func(t *testing.T) {
		pm := &ProductModel{dbpool}

		err := pm.Delete(&Product{})
		if err == nil {
			t.Fatal("expected error")
		}

		if errors.Is(err, ErrRecordNotFound) {
			return
		}

		unexpectedErr(t, err)
	})

	t.Run("create and delete product", func(t *testing.T) {
		pr := &Product{
			Name:        "createAndDeleteTEST",
			Description: "create and delete product",
		}

		pm := &ProductModel{dbpool}

		err := pm.Create(pr)
		unexpectedErr(t, err)

		err = pm.Delete(pr)
		unexpectedErr(t, err)
	})
}
