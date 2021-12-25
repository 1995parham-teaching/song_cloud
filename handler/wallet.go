package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cng-by-example/song_cloud/request"
	"github.com/labstack/echo/v4"
)

type Wallet struct {
	Store *sql.DB
}

// nolint: wrapcheck
func (w *Wallet) Update(c echo.Context) error {
	ctx := c.Request().Context()

	var body request.Wallet
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	stmt, err := w.Store.PrepareContext(
		ctx,
		"UPDATE wallet SET credit = credit + $1 WHERE username = $2",
	)
	if err != nil {
		log.Printf("stmt preparation failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, body.Credit, body.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	n, err := result.RowsAffected()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if n == 0 {
		return echo.ErrNotFound
	}

	return c.NoContent(http.StatusOK)
}

// nolint: wrapcheck, funlen, cyclop
func (w *Wallet) Transfer(c echo.Context) error {
	ctx := c.Request().Context()

	var body request.Transfer
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := w.Store.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	{
		stmt, err := w.Store.PrepareContext(
			ctx,
			"UPDATE wallet SET credit = credit - $1 WHERE username = $2",
		)
		if err != nil {
			log.Printf("stmt preparation failed %s", err)
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		result, err := stmt.ExecContext(ctx, body.Credit, body.Username)
		if err != nil {
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		n, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if n == 0 {
			_ = tx.Rollback()

			return echo.ErrNotFound
		}
	}

	{
		stmt, err := w.Store.PrepareContext(
			ctx,
			"UPDATE wallet SET credit = credit + $1 WHERE username = $2",
		)
		if err != nil {
			log.Printf("stmt preparation failed %s", err)
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		result, err := stmt.ExecContext(ctx, body.Credit, body.EndUser)
		if err != nil {
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		n, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if n == 0 {
			_ = tx.Rollback()

			return echo.ErrNotFound
		}
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (w *Wallet) Register(g *echo.Group) {
	g.POST("/wallet", w.Update)
	g.POST("/transfer", w.Transfer)
}
