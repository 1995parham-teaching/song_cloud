package response

import (
	"time"
)

type User struct {
	Username       string     `json:"username"`
	Password       string     `json:"password"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Email          string     `json:"email"`
	Score          int        `json:"score"`
	SpecialTill    *time.Time `json:"special_till"`
}