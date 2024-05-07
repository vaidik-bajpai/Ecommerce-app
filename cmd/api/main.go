package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// New import
	// New import
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/vaidik-bajpai/ecommerce-api/internal/data"
)

type config struct {
	env  string
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var cfg config

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", 4000, "http network address")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("ECOMMERCE_DB_DSN"), "dsn for the database")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := cfg.openDB()
	if err != nil {
		logger.Fatal(err)
		return
	}

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.port),
		Handler:     app.routes(),
		IdleTimeout: 10 * time.Second,
		ReadTimeout: 30 * time.Second,
	}

	srv.ListenAndServe()
}

func (c config) openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
