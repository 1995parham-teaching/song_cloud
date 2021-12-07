package main

import (
	"errors"
	"log"
	"net/http"


	"song_cloud/config"
	"song_cloud/db"
	"song_cloud/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Read()
	database, err := db.New(cfg.Database)
	if err != nil {
		panic(err)
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

	special := handler.Special{
		Store: database,
	}
	special.Register(app.Group("/api"))

	admin := handler.Admin{
		Store: database,
	}
	admin.Register(app.Group("/api"))

	film := handler.Film{
		Store: database,
	}
	film.Register(app.Group("/api"))

	vote := handler.Vote{
		Store: database,
	}
	vote.Register(app.Group("/api"))

	introduction := handler.Introduction{
		Store: database,
	}
	introduction.Register(app.Group("/api"))

	favorite := handler.Favorite{
		Store: database,
	}
	favorite.Register(app.Group("/api"))

	follow := handler.Follow{
		Store: database,
	}
	follow.Register(app.Group("/api"))

	if err = app.Start(":8080"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("echo initiation failed", err)
	}
}
