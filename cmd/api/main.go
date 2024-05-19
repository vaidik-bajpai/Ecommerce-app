package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"

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
	jwt struct {
		secret string
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

	flag.StringVar(&cfg.jwt.secret, "jwt", os.Getenv("JWT_SECRET"), "secret for json web token")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := cfg.openDB()
	if err != nil {
		logger.Fatal(err)
		return
	}

	defer db.Prisma.Disconnect()

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

func (c config) openDB() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}

	return client, nil

}
