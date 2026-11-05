package request

import (
	"errors"
	"strings"
)

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrNameRequired     = errors.New("playlist name is required")
	ErrNameTooLong      = errors.New("playlist name must be at most 255 characters")
	ErrSongRequired     = errors.New("song is required")
	ErrShareWithMissing = errors.New("user to share with is required")
)

// MaxPlaylistNameLength is the size of the playlist name column.
const MaxPlaylistNameLength = 255

// NewPlaylist creates a playlist for the given user.
type NewPlaylist struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

func (r NewPlaylist) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrUsernameRequired
	}

	name := strings.TrimSpace(r.Name)

	if name == "" {
		return ErrNameRequired
	}

	if len([]rune(name)) > MaxPlaylistNameLength {
		return ErrNameTooLong
	}

	return nil
}

// ListPlaylists lists the playlists that the user owns or are shared with them.
type ListPlaylists struct {
	Username string `query:"username"`
}

func (r ListPlaylists) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrUsernameRequired
	}

	return nil
}

// GetPlaylist reads a playlist on behalf of the given user.
type GetPlaylist struct {
	ID       int    `param:"id"`
	Username string `query:"username"`
}

func (r GetPlaylist) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrUsernameRequired
	}

	return nil
}

// AddPlaylistSong adds a song to a playlist on behalf of its owner.
type AddPlaylistSong struct {
	ID       int    `json:"-"        param:"id"`
	Username string `json:"username"`
	Song     int    `json:"song"`
}

func (r AddPlaylistSong) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrUsernameRequired
	}

	if r.Song <= 0 {
		return ErrSongRequired
	}

	return nil
}

// RemovePlaylistSong removes a song from a playlist on behalf of its owner.
type RemovePlaylistSong struct {
	ID       int    `param:"id"`
	Song     int    `param:"song"`
	Username string `query:"username"`
}

func (r RemovePlaylistSong) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrUsernameRequired
	}

	if r.Song <= 0 {
		return ErrSongRequired
	}

	return nil
}

// SharePlaylist shares a playlist with another user on behalf of its owner.
type SharePlaylist struct {
	ID       int    `json:"-"        param:"id"`
	Username string `json:"username"`
	With     string `json:"with"`
}

func (r SharePlaylist) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrUsernameRequired
	}

	if strings.TrimSpace(r.With) == "" {
		return ErrShareWithMissing
	}

	return nil
}
