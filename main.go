package main

import (
	"song_cloud/config"
	"song_cloud/db"
)

func main() {
	cfg := config.Read()
	database, err := db.New(cfg.Database)
	if err != nil {
		panic(err)
	}


}
