package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cng-by-example/song_cloud/request"
	"github.com/cng-by-example/song_cloud/response"
	"github.com/labstack/echo/v4"
)

type SignUp struct {
	Store *sql.DB
}

// nolint: wrapcheck
func (s *SignUp) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.Signup
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := s.Store.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	{
		stmt, err := s.Store.PrepareContext(
			ctx,
			"INSERT INTO users (username, password, first_name, last_name, email) VALUES ($1, $2, $3, $4, $5)",
		)
		if err != nil {
			log.Printf("stmt preparation failed %s", err)
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		if _, err := stmt.ExecContext(ctx, rq.Username, rq.Password, rq.FirstName, rq.LastName, rq.Email); err != nil {
			log.Printf("stmt exec failed %s", err)
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	{
		stmt, err := s.Store.PrepareContext(ctx, "INSERT INTO wallet (username) VALUES ($1)")
		if err != nil {
			log.Printf("stmt preparation failed %s", err)
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		if _, err := stmt.ExecContext(ctx, rq.Username); err != nil {
			log.Printf("stmt exec failed %s", err)
			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	if rq.Introducer != nil {
		stmt, err := s.Store.PrepareContext(ctx, "INSERT INTO introduce (username, introducer) VALUES ($1, $2)")
		if err != nil {
			log.Printf("stmt preparation failed %s", err)

			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		if _, err := stmt.ExecContext(ctx, rq.Username, *rq.Introducer); err != nil {
			log.Printf("stmt exec failed %s", err)

			_ = tx.Rollback()

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

// Retrieve retrieves URL for given short URL and redirect to it.
// nolint: wrapcheck
func (s *SignUp) Retrieve(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.Login
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user response.User

	stmt, err := s.Store.PrepareContext(ctx, "SELECT * FROM users WHERE username = $1 AND password = $2")
	if err != nil {
		log.Printf("stmt preparation failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, rq.Username, rq.Password).Scan(
		&user.Username, &user.Password, &user.FirstName, &user.LastName,
		&user.Email, &user.SpecialTill, &user.Score); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.ErrUnauthorized
		}

		log.Printf("stmt exec failed %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// nolint: wrapcheck
func (s *SignUp) Update(c echo.Context) error {
	var rq request.Signup
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	query := "UPDATE users SET "

	columns := make(map[string]string)

	if rq.Password != "" {
		columns["password"] = rq.Password
	}

	if rq.FirstName != "" {
		columns["first_name"] = rq.FirstName
	}

	if rq.LastName != "" {
		columns["last_name"] = rq.LastName
	}

	if rq.Email != "" {
		columns["email"] = rq.Email
	}

	for k, v := range columns {
		query += k + " = '" + v + "', "
	}

	query = strings.Trim(query, ", ")

	query += fmt.Sprintf(" WHERE username = '%s'", rq.Username)

	if _, err := s.Store.Exec(query); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (s *SignUp) Register(g *echo.Group) {
	g.POST("/login", s.Retrieve)
	g.POST("/signup", s.Create)
	g.POST("/edit", s.Update)
}
