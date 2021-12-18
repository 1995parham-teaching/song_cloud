package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/elahe-dastan/song_cloud/request"
	"github.com/elahe-dastan/song_cloud/response"

	"github.com/labstack/echo/v4"
)

type SignUp struct {
	Store *sql.DB
}

// todo password constraint doesn't work
// todo unique constraint on email doesn't work
func (s *SignUp) Create(c echo.Context) error {
	var rq request.Signup
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := s.Store.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	query := fmt.Sprintf("INSERT INTO users (username, password, first_name, last_name, email) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		rq.Username, rq.Password, rq.FirstName, rq.LastName, rq.Email)
	if _, err = s.Store.Exec(query); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	query = fmt.Sprintf("INSERT INTO wallet (username) VALUES ('%s')", rq.Username)
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

// Retrieve retrieves URL for given short URL and redirect to it.
// nolint: wrapcheck
func (s *SignUp) Retrieve(c echo.Context) error {
	var rq request.Login
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user response.User
	query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s' AND password = '%s'", rq.Username, rq.Password)
	err := s.Store.QueryRow(query).Scan(&user.Username, &user.Password, &user.FirstName, &user.LastName,
		&user.Email, &user.SpecialTill, &user.Score)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

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

	_, err := s.Store.Exec(query)
	if err != nil {
		return err
	}

	// TODO
	//if result.RowsAffected == 0 {
	//	return ctx.JSON(http.StatusNotFound, DriverSignupError{Message: "referrer not found"})
	//}

	return c.NoContent(http.StatusOK)
}

// Register registers the routes of URL handler on given group.
func (s *SignUp) Register(g *echo.Group) {
	g.POST("/login", s.Retrieve)
	g.POST("/signup", s.Create)
	g.POST("/edit", s.Update)
}
