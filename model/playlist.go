package model

import "time"

// Playlist is a named list of songs that belongs to a single user.
type Playlist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

// PlaylistDetail is a playlist together with its songs.
type PlaylistDetail struct {
	Playlist

	Songs []Song `json:"songs"`
}

// Song is the summary of a song as it is shown inside a playlist.
type Song struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	File  string `json:"file"`
	Price int    `json:"price"`
}
