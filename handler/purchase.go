package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/elahe-dastan/song_cloud/request"
	"github.com/labstack/echo/v4"
)

type Purchase struct {
	Store *sql.DB
}

func (p Purchase) Buy(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.Buy
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	stmt, err := p.Store.PrepareContext(
		ctx,
		"INSERT INTO purchase (username,song_id) VALUES ($1, $2)",
	)
	if err != nil {
		log.Printf("stmt preparation failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, rq.Username, rq.Song); err != nil {
		log.Printf("stmt exec failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (p Purchase) Register(g *echo.Group) {
	g.POST("/buy", p.Buy)
}
