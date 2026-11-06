package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/1995parham-teaching/song_cloud/handler"
	"github.com/1995parham-teaching/song_cloud/model"
	"github.com/1995parham-teaching/song_cloud/store"
	"github.com/labstack/echo/v4"
)

const (
	owner        = "elahe"
	friend       = "parham"
	playlistPath = "/api/playlist"
	songPath     = "/api/playlist/1/song"
	sharePath    = "/api/playlist/1/share"
)

var errDatabaseDown = errors.New("connection refused")

// stubStore records the arguments it receives and answers with canned values.
type stubStore struct {
	playlist  model.Playlist
	detail    model.PlaylistDetail
	playlists []model.Playlist
	err       error

	called   bool
	id       int
	username string
	arg      string
	song     int
}

func (s *stubStore) Create(_ context.Context, owner, name string) (model.Playlist, error) {
	s.called, s.username, s.arg = true, owner, name

	return s.playlist, s.err
}

func (s *stubStore) Get(_ context.Context, id int, viewer string) (model.PlaylistDetail, error) {
	s.called, s.id, s.username = true, id, viewer

	return s.detail, s.err
}

func (s *stubStore) List(_ context.Context, username string) ([]model.Playlist, error) {
	s.called, s.username = true, username

	return s.playlists, s.err
}

func (s *stubStore) AddSong(_ context.Context, id int, username string, songID int) error {
	s.called, s.id, s.username, s.song = true, id, username, songID

	return s.err
}

func (s *stubStore) RemoveSong(_ context.Context, id int, username string, songID int) error {
	s.called, s.id, s.username, s.song = true, id, username, songID

	return s.err
}

func (s *stubStore) Share(_ context.Context, id int, username, with string) error {
	s.called, s.id, s.username, s.arg = true, id, username, with

	return s.err
}

func serve(t *testing.T, s *stubStore, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	app := echo.New()
	handler.Playlist{Store: s}.Register(app.Group("/api"))

	req := httptest.NewRequestWithContext(t.Context(), method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()

	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("cannot decode response %q: %s", rec.Body.String(), err)
	}

	return v
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, want, rec.Body.String())
	}
}

func TestPlaylistCreate(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)
	s := new(stubStore)
	s.playlist = model.Playlist{ID: 7, Name: "road trip", Owner: owner, CreatedAt: createdAt}

	rec := serve(t, s, http.MethodPost, playlistPath, `{"username": "elahe", "name": "  road trip  "}`)

	assertStatus(t, rec, http.StatusCreated)

	if s.username != owner || s.arg != "road trip" {
		t.Errorf("store called with (%q, %q), want (elahe, road trip)", s.username, s.arg)
	}

	if got := decode[model.Playlist](t, rec); got != s.playlist {
		t.Errorf("response = %+v, want %+v", got, s.playlist)
	}
}

func TestPlaylistList(t *testing.T) {
	t.Parallel()

	s := new(stubStore)
	s.playlists = []model.Playlist{
		{ID: 1, Name: "mine", Owner: owner, CreatedAt: time.Time{}},
		{ID: 2, Name: "shared", Owner: friend, CreatedAt: time.Time{}},
	}

	rec := serve(t, s, http.MethodGet, "/api/playlist?username=elahe", "")

	assertStatus(t, rec, http.StatusOK)

	if s.username != owner {
		t.Errorf("store called with %q, want elahe", s.username)
	}

	if got := decode[[]model.Playlist](t, rec); len(got) != len(s.playlists) {
		t.Errorf("response has %d playlists, want %d", len(got), len(s.playlists))
	}
}

func TestPlaylistGet(t *testing.T) {
	t.Parallel()

	s := new(stubStore)
	s.detail = model.PlaylistDetail{
		Playlist: model.Playlist{ID: 3, Name: "focus", Owner: friend, CreatedAt: time.Time{}},
		Songs:    []model.Song{{ID: 1, Name: "first", File: "first.mp3", Price: 0}},
	}

	rec := serve(t, s, http.MethodGet, "/api/playlist/3?username=elahe", "")

	assertStatus(t, rec, http.StatusOK)

	if s.id != 3 || s.username != owner {
		t.Errorf("store called with (%d, %q), want (3, elahe)", s.id, s.username)
	}

	got := decode[model.PlaylistDetail](t, rec)
	if got.ID != 3 || len(got.Songs) != 1 || got.Songs[0] != s.detail.Songs[0] {
		t.Errorf("response = %+v, want %+v", got, s.detail)
	}
}

func TestPlaylistAddSong(t *testing.T) {
	t.Parallel()

	s := new(stubStore)

	// the playlist id comes from the path and cannot be overridden by the body.
	rec := serve(t, s, http.MethodPost, "/api/playlist/5/song", `{"id": 99, "username": "elahe", "song": 2}`)

	assertStatus(t, rec, http.StatusNoContent)

	if s.id != 5 || s.username != owner || s.song != 2 {
		t.Errorf("store called with (%d, %q, %d), want (5, elahe, 2)", s.id, s.username, s.song)
	}
}

func TestPlaylistRemoveSong(t *testing.T) {
	t.Parallel()

	s := new(stubStore)

	rec := serve(t, s, http.MethodDelete, "/api/playlist/5/song/2?username=elahe", "")

	assertStatus(t, rec, http.StatusNoContent)

	if s.id != 5 || s.username != owner || s.song != 2 {
		t.Errorf("store called with (%d, %q, %d), want (5, elahe, 2)", s.id, s.username, s.song)
	}
}

func TestPlaylistShare(t *testing.T) {
	t.Parallel()

	s := new(stubStore)

	rec := serve(t, s, http.MethodPost, "/api/playlist/5/share", `{"username": "elahe", "with": "parham"}`)

	assertStatus(t, rec, http.StatusNoContent)

	if s.id != 5 || s.username != owner || s.arg != friend {
		t.Errorf("store called with (%d, %q, %q), want (5, elahe, parham)", s.id, s.username, s.arg)
	}
}

func TestPlaylistValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		target string
		body   string
	}{
		{"create without username", http.MethodPost, playlistPath, `{"name": "road trip"}`},
		{"create without name", http.MethodPost, playlistPath, `{"username": "elahe", "name": "   "}`},
		{
			"create with too long name", http.MethodPost, playlistPath,
			fmt.Sprintf(`{"username": "elahe", "name": %q}`, strings.Repeat("a", 256)),
		},
		{"create with malformed body", http.MethodPost, playlistPath, `{"username": `},
		{"list without username", http.MethodGet, playlistPath, ""},
		{"get without username", http.MethodGet, "/api/playlist/1", ""},
		{"get with non numeric id", http.MethodGet, "/api/playlist/abc?username=elahe", ""},
		{"add song without song", http.MethodPost, songPath, `{"username": "elahe"}`},
		{"add song without username", http.MethodPost, songPath, `{"song": 1}`},
		{"remove song without username", http.MethodDelete, songPath + "/2", ""},
		{"remove song with non numeric song", http.MethodDelete, songPath + "/abc?username=elahe", ""},
		{"remove song with invalid song", http.MethodDelete, songPath + "/0?username=elahe", ""},
		{"share without target", http.MethodPost, sharePath, `{"username": "elahe"}`},
		{"share without username", http.MethodPost, sharePath, `{"with": "parham"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := new(stubStore)

			rec := serve(t, s, tc.method, tc.target, tc.body)

			assertStatus(t, rec, http.StatusBadRequest)

			if s.called {
				t.Error("store must not be called for an invalid request")
			}
		})
	}
}

func TestPlaylistStoreErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		status int
	}{
		{"not found", fmt.Errorf("playlist 1: %w", store.ErrNotFound), http.StatusNotFound},
		{"forbidden", fmt.Errorf("playlist 1: %w", store.ErrForbidden), http.StatusForbidden},
		{"duplicate", fmt.Errorf("playlist: %w", store.ErrDuplicate), http.StatusConflict},
		{"invalid", fmt.Errorf("self share: %w", store.ErrInvalid), http.StatusBadRequest},
		{"internal", errDatabaseDown, http.StatusInternalServerError},
	}

	requests := []struct {
		method string
		target string
		body   string
	}{
		{http.MethodPost, playlistPath, `{"username": "elahe", "name": "road trip"}`},
		{http.MethodGet, "/api/playlist?username=elahe", ""},
		{http.MethodGet, "/api/playlist/1?username=elahe", ""},
		{http.MethodPost, songPath, `{"username": "elahe", "song": 1}`},
		{http.MethodDelete, songPath + "/1?username=elahe", ""},
		{http.MethodPost, sharePath, `{"username": "elahe", "with": "parham"}`},
	}

	for _, tc := range tests {
		for _, rq := range requests {
			t.Run(fmt.Sprintf("%s %s %s", tc.name, rq.method, rq.target), func(t *testing.T) {
				t.Parallel()

				s := new(stubStore)
				s.err = tc.err

				rec := serve(t, s, rq.method, rq.target, rq.body)

				assertStatus(t, rec, tc.status)

				if strings.Contains(rec.Body.String(), errDatabaseDown.Error()) {
					t.Error("internal errors must not be leaked to the client")
				}
			})
		}
	}
}
