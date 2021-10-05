package main

import (
	"context"
	"log"
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
		log.Fatal("ERROR: no DATABASE_URL\n")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("ERROR: no PORT")
		return
	}

	dbpool, err := ConnectDb(connstr)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
		return
	}

	app.models = db.NewModels(dbpool)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	log.Printf("listening on port %s\n", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
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
