package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/elahe-dastan/song_cloud/config"
	"github.com/elahe-dastan/song_cloud/db"
	"github.com/elahe-dastan/song_cloud/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	app := echo.New()
	app.Use(middleware.CORS())

	signup := handler.SignUp{
		Store: database,
	}
	signup.Register(app.Group("/api"))

	wallet := handler.Wallet{
		Store: database,
	}
	wallet.Register(app.Group("/api"))

	premium := handler.Premium{
		Store: database,
	}
	premium.Register(app.Group("/api"))

	song := handler.Song{
		Store: database,
	}
	song.Register(app.Group("/api"))

	purchase := handler.Purchase{
		Store: database,
	}
	purchase.Register(app.Group("/api"))

	if err = app.Start(":8080"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("echo initiation failed", err)
	}
}

// Register server command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "serve",
			Short: "Run server to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
