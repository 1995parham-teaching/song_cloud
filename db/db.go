package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/1995parham-teaching/song_cloud/config"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// New creates a new postgres connection and tests it.
func New(cfg config.Database) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLmode)

	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Println("Successfully connected!")

	return db, nil
}
