package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	models db.Models
}

func main() {
	app := &application{}
	connstr := os.Getenv("DATABASE_URL")
	dbpool, err := ConnectDb(connstr)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	app.models = db.NewModels(dbpool)
	port := flag.Int("port", 8080, "port to listen on")

	srv := http.Server{
		Addr:    ":" + strconv.Itoa(*port),
		Handler: app.routes(),
	}

	fmt.Printf("listening on port %d\n", *port)

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
}

func ConnectDb(connstr string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
