package store_test

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/1995parham-teaching/song_cloud/config"
	"github.com/1995parham-teaching/song_cloud/db"
	"github.com/1995parham-teaching/song_cloud/model"
	"github.com/1995parham-teaching/song_cloud/store"
)

// fixture creates isolated records for a test and removes them when it finishes.
type fixture struct {
	db    *sql.DB
	store store.Playlist
}

// connect opens the configured database. These tests need a migrated database,
// so they are skipped locally when it is not reachable and fail on CI.
func connect(t *testing.T) fixture {
	t.Helper()

	database, err := db.New(config.Read().Database)
	if err != nil {
		if os.Getenv("CI") != "" {
			t.Fatalf("database is required on CI: %s", err)
		}

		t.Skipf("database is not available: %s", err)
	}

	t.Cleanup(func() { _ = database.Close() })

	return fixture{db: database, store: store.Playlist{DB: database}}
}

func (f fixture) exec(t *testing.T, query string, args ...any) {
	t.Helper()

	if _, err := f.db.ExecContext(context.WithoutCancel(t.Context()), query, args...); err != nil {
		t.Fatalf("%s: %s", query, err)
	}
}

func randomName(prefix string) string {
	return prefix + "_" + strings.ToLower(rand.Text())
}

// user creates a user with a random name. Deleting it cascades to its playlists and shares.
func (f fixture) user(t *testing.T) string {
	t.Helper()

	username := randomName("test")

	f.exec(t,
		"INSERT INTO users (username, password, email) VALUES ($1::text, 'password1', $1::text || '@example.com')",
		username,
	)
	t.Cleanup(func() { f.exec(t, "DELETE FROM users WHERE username = $1", username) })

	return username
}

func (f fixture) song(t *testing.T, name string) int {
	t.Helper()

	var id int

	if err := f.db.QueryRowContext(
		t.Context(), "INSERT INTO song (name, file, price) VALUES ($1::text, $1::text || '.mp3', 0) RETURNING id", name,
	).Scan(&id); err != nil {
		t.Fatalf("insert song: %s", err)
	}

	t.Cleanup(func() { f.exec(t, "DELETE FROM song WHERE id = $1", id) })

	return id
}

func (f fixture) playlist(t *testing.T, owner string) model.Playlist {
	t.Helper()

	playlist, err := f.store.Create(t.Context(), owner, "playlist of "+owner)
	if err != nil {
		t.Fatalf("create playlist: %s", err)
	}

	return playlist
}

func (f fixture) addSong(t *testing.T, playlist int, owner string, song int) {
	t.Helper()

	if err := f.store.AddSong(t.Context(), playlist, owner, song); err != nil {
		t.Fatalf("add song %d to playlist %d: %s", song, playlist, err)
	}
}

func (f fixture) share(t *testing.T, playlist int, owner, with string) {
	t.Helper()

	if err := f.store.Share(t.Context(), playlist, owner, with); err != nil {
		t.Fatalf("share playlist %d with %q: %s", playlist, with, err)
	}
}

func (f fixture) get(t *testing.T, playlist int, viewer string) model.PlaylistDetail {
	t.Helper()

	detail, err := f.store.Get(t.Context(), playlist, viewer)
	if err != nil {
		t.Fatalf("get playlist %d as %q: %s", playlist, viewer, err)
	}

	return detail
}

func (f fixture) list(t *testing.T, username string) []model.Playlist {
	t.Helper()

	playlists, err := f.store.List(t.Context(), username)
	if err != nil {
		t.Fatalf("list playlists of %q: %s", username, err)
	}

	return playlists
}

func assertErr(t *testing.T, err, want error) {
	t.Helper()

	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func songIDs(songs []model.Song) []int {
	ids := make([]int, 0, len(songs))
	for _, s := range songs {
		ids = append(ids, s.ID)
	}

	return ids
}

func playlistIDs(playlists []model.Playlist) []int {
	ids := make([]int, 0, len(playlists))
	for _, p := range playlists {
		ids = append(ids, p.ID)
	}

	return ids
}

func assertIDs(t *testing.T, got []int, want ...int) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("ids = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ids = %v, want %v", got, want)
		}
	}
}

// TestPlaylist runs each scenario as a parallel subtest with its own users, songs and playlists.
//
// nolint: paralleltest // every scenario calls t.Parallel, the linter cannot follow method values.
func TestPlaylist(t *testing.T) {
	t.Parallel()

	f := connect(t)

	t.Run("create", f.testCreate)
	t.Run("create for unknown user", f.testCreateUnknownUser)
	t.Run("add songs", f.testAddSongs)
	t.Run("add song errors", f.testAddSongErrors)
	t.Run("remove song", f.testRemoveSong)
	t.Run("empty playlist has no songs", f.testEmptyPlaylist)
	t.Run("get unknown playlist", f.testGetUnknown)
	t.Run("share", f.testShare)
	t.Run("share errors", f.testShareErrors)
	t.Run("list", f.testList)
}

func (f fixture) testCreate(t *testing.T) {
	t.Parallel()

	owner := f.user(t)

	playlist, err := f.store.Create(t.Context(), owner, "road trip")
	if err != nil {
		t.Fatalf("create: %s", err)
	}

	if playlist.ID == 0 || playlist.Name != "road trip" || playlist.Owner != owner || playlist.CreatedAt.IsZero() {
		t.Errorf("unexpected playlist %+v", playlist)
	}

	_, err = f.store.Create(t.Context(), owner, "road trip")
	assertErr(t, err, store.ErrDuplicate)

	// the same name is fine for another user.
	if _, err := f.store.Create(t.Context(), f.user(t), "road trip"); err != nil {
		t.Errorf("create with the same name for another user: %s", err)
	}
}

func (f fixture) testCreateUnknownUser(t *testing.T) {
	t.Parallel()

	_, err := f.store.Create(t.Context(), randomName("unknown"), "road trip")
	assertErr(t, err, store.ErrNotFound)
}

func (f fixture) testAddSongs(t *testing.T) {
	t.Parallel()

	owner := f.user(t)
	playlist := f.playlist(t, owner)
	first, second := f.song(t, "first"), f.song(t, "second")

	// adding a song twice is a no-op.
	f.addSong(t, playlist.ID, owner, first)
	f.addSong(t, playlist.ID, owner, second)
	f.addSong(t, playlist.ID, owner, first)

	detail := f.get(t, playlist.ID, owner)

	if detail.Playlist != playlist {
		t.Errorf("playlist = %+v, want %+v", detail.Playlist, playlist)
	}

	assertIDs(t, songIDs(detail.Songs), first, second)

	if want := (model.Song{ID: first, Name: "first", File: "first.mp3", Price: 0}); detail.Songs[0] != want {
		t.Errorf("song = %+v, want %+v", detail.Songs[0], want)
	}
}

func (f fixture) testAddSongErrors(t *testing.T) {
	t.Parallel()

	owner := f.user(t)
	playlist := f.playlist(t, owner)
	song := f.song(t, "song")

	assertErr(t, f.store.AddSong(t.Context(), playlist.ID, f.user(t), song), store.ErrForbidden)
	assertErr(t, f.store.AddSong(t.Context(), playlist.ID, owner, -1), store.ErrNotFound)
	assertErr(t, f.store.AddSong(t.Context(), -1, owner, song), store.ErrNotFound)
}

func (f fixture) testRemoveSong(t *testing.T) {
	t.Parallel()

	owner, friend := f.user(t), f.user(t)
	playlist := f.playlist(t, owner)
	first, second := f.song(t, "first"), f.song(t, "second")

	f.addSong(t, playlist.ID, owner, first)
	f.addSong(t, playlist.ID, owner, second)
	f.share(t, playlist.ID, owner, friend)

	// shared users have read-only access.
	assertErr(t, f.store.RemoveSong(t.Context(), playlist.ID, friend, first), store.ErrForbidden)

	if err := f.store.RemoveSong(t.Context(), playlist.ID, owner, first); err != nil {
		t.Fatalf("remove song: %s", err)
	}

	assertIDs(t, songIDs(f.get(t, playlist.ID, owner).Songs), second)

	// the song is no longer in the playlist.
	assertErr(t, f.store.RemoveSong(t.Context(), playlist.ID, owner, first), store.ErrNotFound)
	assertErr(t, f.store.RemoveSong(t.Context(), -1, owner, second), store.ErrNotFound)

	// the song itself is not deleted and can be added again.
	f.addSong(t, playlist.ID, owner, first)
	assertIDs(t, songIDs(f.get(t, playlist.ID, owner).Songs), second, first)
}

func (f fixture) testEmptyPlaylist(t *testing.T) {
	t.Parallel()

	owner := f.user(t)
	playlist := f.playlist(t, owner)

	// an empty slice (not nil) is encoded as [] instead of null.
	if songs := f.get(t, playlist.ID, owner).Songs; songs == nil || len(songs) != 0 {
		t.Errorf("songs = %#v, want an empty slice", songs)
	}
}

func (f fixture) testGetUnknown(t *testing.T) {
	t.Parallel()

	_, err := f.store.Get(t.Context(), -1, f.user(t))
	assertErr(t, err, store.ErrNotFound)
}

func (f fixture) testShare(t *testing.T) {
	t.Parallel()

	owner, friend, stranger := f.user(t), f.user(t), f.user(t)
	playlist := f.playlist(t, owner)
	song := f.song(t, "shared song")

	f.addSong(t, playlist.ID, owner, song)

	_, err := f.store.Get(t.Context(), playlist.ID, friend)
	assertErr(t, err, store.ErrForbidden)

	// sharing twice is a no-op.
	f.share(t, playlist.ID, owner, friend)
	f.share(t, playlist.ID, owner, friend)

	assertIDs(t, songIDs(f.get(t, playlist.ID, friend).Songs), song)

	_, err = f.store.Get(t.Context(), playlist.ID, stranger)
	assertErr(t, err, store.ErrForbidden)

	// shared users have read-only access.
	assertErr(t, f.store.AddSong(t.Context(), playlist.ID, friend, song), store.ErrForbidden)
	assertErr(t, f.store.Share(t.Context(), playlist.ID, friend, stranger), store.ErrForbidden)
}

func (f fixture) testShareErrors(t *testing.T) {
	t.Parallel()

	owner := f.user(t)
	playlist := f.playlist(t, owner)

	assertErr(t, f.store.Share(t.Context(), playlist.ID, owner, owner), store.ErrInvalid)
	assertErr(t, f.store.Share(t.Context(), playlist.ID, owner, randomName("unknown")), store.ErrNotFound)
	assertErr(t, f.store.Share(t.Context(), -1, owner, f.user(t)), store.ErrNotFound)
}

func (f fixture) testList(t *testing.T) {
	t.Parallel()

	owner, friend, stranger := f.user(t), f.user(t), f.user(t)
	owned := f.playlist(t, owner)
	shared := f.playlist(t, friend)

	f.share(t, shared.ID, friend, owner)

	// playlists are ordered by id, so the owned one (created first) comes first.
	assertIDs(t, playlistIDs(f.list(t, owner)), owned.ID, shared.ID)
	assertIDs(t, playlistIDs(f.list(t, friend)), shared.ID)

	// an empty slice (not nil) is encoded as [] instead of null.
	if playlists := f.list(t, stranger); playlists == nil || len(playlists) != 0 {
		t.Errorf("playlists = %#v, want an empty slice", playlists)
	}
}
