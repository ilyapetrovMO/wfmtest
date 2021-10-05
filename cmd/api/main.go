package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

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
		log.Fatalf("ERROR: %s\n%s", err, connstr)
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
	// Give time for postgres to finish its setup script(in docker-compose)
	time.Sleep(time.Second * 2)
	dbpool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err = dbpool.Ping(ctx)

	return dbpool, err
}
