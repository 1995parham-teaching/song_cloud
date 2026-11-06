package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/1995parham-teaching/song_cloud/model"
	"github.com/1995parham-teaching/song_cloud/request"
	"github.com/1995parham-teaching/song_cloud/store"
	"github.com/labstack/echo/v4"
)

// PlaylistStore is the persistence layer that the playlist handler depends on.
type PlaylistStore interface {
	Create(ctx context.Context, owner, name string) (model.Playlist, error)
	Get(ctx context.Context, id int, viewer string) (model.PlaylistDetail, error)
	List(ctx context.Context, username string) ([]model.Playlist, error)
	AddSong(ctx context.Context, id int, username string, songID int) error
	RemoveSong(ctx context.Context, id int, username string, songID int) error
	Share(ctx context.Context, id int, username, with string) error
}

type Playlist struct {
	Store PlaylistStore
}

// Create creates a new empty playlist.
func (p Playlist) Create(c echo.Context) error {
	rq, err := bind[request.NewPlaylist](c)
	if err != nil {
		return err
	}

	playlist, err := p.Store.Create(c.Request().Context(), rq.Username, strings.TrimSpace(rq.Name))
	if err != nil {
		return playlistError(err)
	}

	return c.JSON(http.StatusCreated, playlist)
}

// List returns the playlists that the user owns or are shared with them.
func (p Playlist) List(c echo.Context) error {
	rq, err := bind[request.ListPlaylists](c)
	if err != nil {
		return err
	}

	playlists, err := p.Store.List(c.Request().Context(), rq.Username)
	if err != nil {
		return playlistError(err)
	}

	return c.JSON(http.StatusOK, playlists)
}

// Get returns a playlist with its songs.
func (p Playlist) Get(c echo.Context) error {
	rq, err := bind[request.GetPlaylist](c)
	if err != nil {
		return err
	}

	playlist, err := p.Store.Get(c.Request().Context(), rq.ID, rq.Username)
	if err != nil {
		return playlistError(err)
	}

	return c.JSON(http.StatusOK, playlist)
}

// AddSong adds a song to the playlist.
func (p Playlist) AddSong(c echo.Context) error {
	rq, err := bind[request.AddPlaylistSong](c)
	if err != nil {
		return err
	}

	if err := p.Store.AddSong(c.Request().Context(), rq.ID, rq.Username, rq.Song); err != nil {
		return playlistError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveSong removes a song from the playlist.
func (p Playlist) RemoveSong(c echo.Context) error {
	rq, err := bind[request.RemovePlaylistSong](c)
	if err != nil {
		return err
	}

	if err := p.Store.RemoveSong(c.Request().Context(), rq.ID, rq.Username, rq.Song); err != nil {
		return playlistError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Share shares the playlist with another user.
func (p Playlist) Share(c echo.Context) error {
	rq, err := bind[request.SharePlaylist](c)
	if err != nil {
		return err
	}

	if err := p.Store.Share(c.Request().Context(), rq.ID, rq.Username, rq.With); err != nil {
		return playlistError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (p Playlist) Register(g *echo.Group) {
	g.POST("/playlist", p.Create)
	g.GET("/playlist", p.List)
	g.GET("/playlist/:id", p.Get)
	g.POST("/playlist/:id/song", p.AddSong)
	g.DELETE("/playlist/:id/song/:song", p.RemoveSong)
	g.POST("/playlist/:id/share", p.Share)
}

// bind binds the request into T and validates it.
func bind[T interface{ Validate() error }](c echo.Context) (T, error) {
	var rq T

	// echo already reports binding failures as *echo.HTTPError (400 or 415).
	if err := c.Bind(&rq); err != nil {
		return rq, err
	}

	if err := rq.Validate(); err != nil {
		return rq, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return rq, nil
}

// playlistError maps store errors into HTTP errors.
func playlistError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrForbidden):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case errors.Is(err, store.ErrDuplicate):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case errors.Is(err, store.ErrInvalid):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		log.Printf("playlist store failed: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
