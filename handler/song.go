package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/elahe-dastan/song_cloud/request"
	"github.com/labstack/echo/v4"
)

type Song struct {
	Store *sql.DB
}

func (s Song) Play(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.PlaySong
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	stmt, err := s.Store.PrepareContext(
		ctx,
		"CALL play ($1, $2)",
	)
	if err != nil {
		log.Printf("stmt preparation failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, rq.ID, rq.Username); err != nil {
		log.Printf("stmt exec failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s Song) New(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.NewSong
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	stmt, err := s.Store.PrepareContext(
		ctx,
		"INSERT INTO song (name, file, production_year, explanation, price) VALUES ($1, $2, $3, $4, $5)",
	)
	if err != nil {
		log.Printf("stmt preparation failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, rq.Name, rq.File, rq.ProductionYear, rq.Explanation, rq.Price); err != nil {
		log.Printf("stmt exec failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (s Song) Register(g *echo.Group) {
	g.POST("/song", s.New)
	g.POST("/play", s.Play)
}
