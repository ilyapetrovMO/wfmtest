package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var connstring = "postgres://postgres:postgres@127.0.0.1:5432/wfmtest"

func TestCreateProductHandler(t *testing.T) {
	pool, err := connectDbTest(connstring)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer pool.Close()
	defer pool.Exec(context.Background(), `DELETE FROM products WHERE name LIKE '%TEST'`)

	app := &application{
		models: db.NewModels(pool),
		logger: logrus.New(),
	}

	t.Run("well formed product", func(t *testing.T) {
		want := `{"name": "testCreateTEST", "description": "desc"}`

		b := strings.NewReader(want)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://example.com/products", b)

		app.createProductHandler(w, r)
		resp := w.Result()

		js := json.NewDecoder(resp.Body)
		got := &struct{ Product db.Product }{}
		js.Decode(got)

		if got.Product.Name != "testCreateTEST" {
			t.Errorf("got %v, want %s", got.Product, want)
		}
	})

	t.Run("malformed product", func(t *testing.T) {
		want := `"name": "must be provided"`
		input := `{"name":""}`

		b := strings.NewReader(input)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://example.com/products", b)

		app.createProductHandler(w, r)
		resp := w.Result()
		got, _ := io.ReadAll(resp.Body)

		if !strings.Contains(string(got), want) {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func connectDbTest(connstr string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err = dbpool.Ping(ctx)

	return dbpool, err
}
