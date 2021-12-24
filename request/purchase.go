package request

type Buy struct {
	Username string `json:"username"`
	Song     int    `json:"song"`
}
