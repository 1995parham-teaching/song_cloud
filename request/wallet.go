package request

type Wallet struct {
	Username string `json:"username"`
	Credit   int    `json:"credit"`
}

type Transfer struct {
	Username string `json:"username"`
	EndUser  string `json:"end_user"`
	Credit   int    `json:"credit"`
}
