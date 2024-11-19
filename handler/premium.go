package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/1995parham-teaching/song_cloud/request"
	"github.com/labstack/echo/v4"
)

type Premium struct {
	Store *sql.DB
}

const HoursInDay = 24

// nolint: wrapcheck
func (p Premium) Extend(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.Premium
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	stmt, err := p.Store.PrepareContext(ctx, "UPDATE users SET premium_till = $1 WHERE username = $2")
	if err != nil {
		log.Printf("stmt preparation failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(
		ctx,
		time.Now().Add(time.Duration(rq.Duration)*time.Hour*HoursInDay),
		rq.Username,
	); err != nil {
		log.Printf("stmt exec failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (p Premium) Register(g *echo.Group) {
	g.POST("/extend", p.Extend)
}
