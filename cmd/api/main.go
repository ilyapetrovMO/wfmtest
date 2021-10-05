package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	models db.Models
}

func main() {
	app := &application{}
	connstr := os.Getenv("DATABASE_URL")
	if connstr == "" {
		fmt.Print("ERROR: no DATABASE_URL\n")
		return
	}
	portstr := os.Getenv("PORT")
	if portstr == "" {
		fmt.Print("ERROR: no PORT\n")
		return
	}

	dbpool, err := ConnectDb(connstr)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	app.models = db.NewModels(dbpool)

	srv := http.Server{
		Addr:    ":" + portstr,
		Handler: app.routes(),
	}

	fmt.Printf("listening on port %s\n", portstr)

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
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
