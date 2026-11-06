package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/1995parham-teaching/song_cloud/model"
)

// Playlist persists playlists in PostgreSQL and enforces who can see or change them:
// only the owner can modify or share a playlist, and the owner and the users it is
// shared with can read it.
type Playlist struct {
	DB *sql.DB
}

// Create creates an empty playlist for the given owner.
func (p Playlist) Create(ctx context.Context, owner, name string) (model.Playlist, error) {
	var playlist model.Playlist

	err := p.DB.QueryRowContext(
		ctx,
		"INSERT INTO playlist (name, owner) VALUES ($1, $2) RETURNING id, name, owner, created_at",
		name, owner,
	).Scan(&playlist.ID, &playlist.Name, &playlist.Owner, &playlist.CreatedAt)

	switch {
	case err == nil:
		return playlist, nil
	case hasCode(err, codeForeignKeyViolation):
		return model.Playlist{}, fmt.Errorf("user %q: %w", owner, ErrNotFound)
	case hasCode(err, codeUniqueViolation):
		return model.Playlist{}, fmt.Errorf("playlist %q of user %q: %w", name, owner, ErrDuplicate)
	default:
		return model.Playlist{}, fmt.Errorf("insert playlist: %w", err)
	}
}

// Get returns a playlist with its songs if the viewer is its owner or it is shared with them.
func (p Playlist) Get(ctx context.Context, id int, viewer string) (model.PlaylistDetail, error) {
	var (
		detail  model.PlaylistDetail
		allowed bool
	)

	err := p.DB.QueryRowContext(ctx, `
		SELECT p.id, p.name, p.owner, p.created_at,
		       p.owner = $2 OR EXISTS (
		           SELECT FROM playlist_share s WHERE s.playlist_id = p.id AND s.username = $2
		       )
		FROM playlist p
		WHERE p.id = $1`,
		id, viewer,
	).Scan(&detail.ID, &detail.Name, &detail.Owner, &detail.CreatedAt, &allowed)
	if errors.Is(err, sql.ErrNoRows) {
		return model.PlaylistDetail{}, fmt.Errorf("playlist %d: %w", id, ErrNotFound)
	}

	if err != nil {
		return model.PlaylistDetail{}, fmt.Errorf("select playlist: %w", err)
	}

	if !allowed {
		return model.PlaylistDetail{}, fmt.Errorf("playlist %d is not shared with %q: %w", id, viewer, ErrForbidden)
	}

	detail.Songs, err = p.songs(ctx, id)
	if err != nil {
		return model.PlaylistDetail{}, err
	}

	return detail, nil
}

// List returns the playlists that the user owns or that are shared with them.
func (p Playlist) List(ctx context.Context, username string) ([]model.Playlist, error) {
	rows, err := p.DB.QueryContext(ctx, `
		SELECT p.id, p.name, p.owner, p.created_at
		FROM playlist p
		WHERE p.owner = $1
		   OR EXISTS (SELECT FROM playlist_share s WHERE s.playlist_id = p.id AND s.username = $1)
		ORDER BY p.id`,
		username,
	)
	if err != nil {
		return nil, fmt.Errorf("select playlists: %w", err)
	}
	defer rows.Close()

	playlists := make([]model.Playlist, 0)

	for rows.Next() {
		var playlist model.Playlist

		if err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.Owner, &playlist.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan playlist: %w", err)
		}

		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlists: %w", err)
	}

	return playlists, nil
}

// AddSong adds a song to the playlist. Only the owner can add songs and adding
// a song that is already in the playlist is a no-op.
func (p Playlist) AddSong(ctx context.Context, id int, username string, songID int) error {
	if err := p.authorizeOwner(ctx, id, username); err != nil {
		return err
	}

	_, err := p.DB.ExecContext(
		ctx,
		"INSERT INTO playlist_song (playlist_id, song_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		id, songID,
	)

	switch {
	case err == nil:
		return nil
	case hasCode(err, codeForeignKeyViolation):
		return fmt.Errorf("song %d: %w", songID, ErrNotFound)
	default:
		return fmt.Errorf("insert playlist song: %w", err)
	}
}

// RemoveSong removes a song from the playlist. Only the owner can remove songs.
func (p Playlist) RemoveSong(ctx context.Context, id int, username string, songID int) error {
	if err := p.authorizeOwner(ctx, id, username); err != nil {
		return err
	}

	result, err := p.DB.ExecContext(
		ctx,
		"DELETE FROM playlist_song WHERE playlist_id = $1 AND song_id = $2",
		id, songID,
	)
	if err != nil {
		return fmt.Errorf("delete playlist song: %w", err)
	}

	removed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete playlist song: %w", err)
	}

	if removed == 0 {
		return fmt.Errorf("song %d is not in playlist %d: %w", songID, id, ErrNotFound)
	}

	return nil
}

// Share gives another user read access to the playlist. Only the owner can share
// a playlist and sharing it again with the same user is a no-op.
func (p Playlist) Share(ctx context.Context, id int, username, with string) error {
	if err := p.authorizeOwner(ctx, id, username); err != nil {
		return err
	}

	if with == username {
		return fmt.Errorf("playlist cannot be shared with its owner: %w", ErrInvalid)
	}

	_, err := p.DB.ExecContext(
		ctx,
		"INSERT INTO playlist_share (playlist_id, username) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		id, with,
	)

	switch {
	case err == nil:
		return nil
	case hasCode(err, codeForeignKeyViolation):
		return fmt.Errorf("user %q: %w", with, ErrNotFound)
	default:
		return fmt.Errorf("insert playlist share: %w", err)
	}
}

// authorizeOwner makes sure the playlist exists and belongs to the given user.
func (p Playlist) authorizeOwner(ctx context.Context, id int, username string) error {
	var owner string

	err := p.DB.QueryRowContext(ctx, "SELECT owner FROM playlist WHERE id = $1", id).Scan(&owner)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("playlist %d: %w", id, ErrNotFound)
	}

	if err != nil {
		return fmt.Errorf("select playlist owner: %w", err)
	}

	if owner != username {
		return fmt.Errorf("playlist %d does not belong to %q: %w", id, username, ErrForbidden)
	}

	return nil
}

// songs returns the songs of a playlist in the order they were added.
func (p Playlist) songs(ctx context.Context, id int) ([]model.Song, error) {
	rows, err := p.DB.QueryContext(ctx, `
		SELECT s.id, COALESCE(s.name, ''), COALESCE(s.file, ''), COALESCE(s.price, 0)
		FROM playlist_song ps
		JOIN song s ON s.id = ps.song_id
		WHERE ps.playlist_id = $1
		ORDER BY ps.added_at, s.id`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("select playlist songs: %w", err)
	}
	defer rows.Close()

	songs := make([]model.Song, 0)

	for rows.Next() {
		var song model.Song

		if err := rows.Scan(&song.ID, &song.Name, &song.File, &song.Price); err != nil {
			return nil, fmt.Errorf("scan playlist song: %w", err)
		}

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playlist songs: %w", err)
	}

	return songs, nil
}
