package request

type NewSong struct {
	Name           string `json:"name"`
	File           string `json:"file"`
	ProductionYear int    `json:"production_year"`
	Explanation    string `json:"explanation"`
	Price          int    `json:"price"`
}

type PlaySong struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type AssignSong struct {
	ID       int `json:"id"`
	Category int `json:"category"`
}
