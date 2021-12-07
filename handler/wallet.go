package handler

import (
	"database/sql"
	"fmt"
	"net/http"


	"github.com/labstack/echo/v4"
	"song_cloud/request"
)

type Wallet struct {
	Store *sql.DB
}

func (w *Wallet) Update(c echo.Context) error {
	var body request.Wallet
	err := c.Bind(&body)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE wallet SET credit = credit + %d WHERE username = '%s'", body.Credit, body.Username)
	_, err = w.Store.Query(query)
	if err != nil {
		return err
	}

	// todo
	//if result.RowsAffected == 0 {
	//	return ctx.JSON(http.StatusNotFound, DriverSignupError{Message: "referrer not found"})
	//}

	return c.NoContent(http.StatusOK)
}

func (w *Wallet) Transfer(c echo.Context) error {
	var body request.Transfer
	err := c.Bind(&body)
	if err != nil {
		return err
	}

	tx, err := w.Store.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	query := fmt.Sprintf("UPDATE wallet SET credit = credit - %d WHERE username = '%s'", body.Credit, body.Username)
	_, err = w.Store.Query(query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	query = fmt.Sprintf("UPDATE wallet SET credit = credit + %d WHERE username = '%s'", body.Credit, body.EndUser)
	_, err = w.Store.Query(query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// todo
	//if result.RowsAffected == 0 {
	//	return ctx.JSON(http.StatusNotFound, DriverSignupError{Message: "referrer not found"})
	//}

	return c.NoContent(http.StatusOK)
}

// Register registers the routes of URL handler on given group.
func (w *Wallet) Register(g *echo.Group) {
	g.POST("/wallet", w.Update)
	g.POST("/transfer", w.Transfer)
}