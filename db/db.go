package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/1995parham-teaching/song_cloud/config"
)

// New creates a new postgres connection and tests it.
func New(cfg config.Database) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	if err = db.PingContext(context.Background()); err != nil {
		panic(err)
	}

	log.Println("Successfully connected!")

	return db, nil
}
