package migrate

import (
	"errors"

	"github.com/1995parham-teaching/song_cloud/config"
	"github.com/1995parham-teaching/song_cloud/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // supports migration with files.
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq" // loads postgresql driver.
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatal("database initiation failed", err)
	}

	// nolint: exhaustruct
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migration",
		"postgres", driver)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Print("database is up to date")

			return
		}

		log.Fatalf(err.Error())
	}
}

// Register migrate command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustruct
		&cobra.Command{
			Use:   "migrate",
			Short: "database migration",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
