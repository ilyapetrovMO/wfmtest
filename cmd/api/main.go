package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type application struct {
	models db.Models
	logger *logrus.Logger
}

func main() {
	log := logrus.New()

	app := &application{
		logger: log,
	}

	connstr := os.Getenv("DATABASE_URL")
	if connstr == "" {
		app.logger.Fatalf("ERROR: no DATABASE_URL\n")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		app.logger.Fatal("ERROR: no PORT")
		return
	}

	dbpool, err := connectDb(connstr)
	if err != nil {
		app.logger.Fatalf("ERROR: %s\n%s", err, connstr)
		return
	}

	app.models = db.NewModels(dbpool)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: app.routes(),
	}

	app.logger.Printf("listening on port %s\n", port)

	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatalf("ERROR: %s\n", err)
		return
	}
}

func connectDb(connstr string) (*pgxpool.Pool, error) {
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
