/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package db

import (
	"simmie/src/utils"
	"crypto/rand"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Song struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	TrackID  string `json:"track_id"`
	Duration int    `json:"duration"`
	Platform string `json:"platform"`
}

type Playlist struct {
	ID     string
	Name   string
	UserID int64
	Songs  []Song
}

func generateUniquePlaylistID() string {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	return fmt.Sprintf("tgpl_%x", b)
}

func (db *Database) CreatePlaylist(name string, userID int64) (string, error) {
	ctx, cancel := db.ctx()
	defer cancel()

	id := generateUniquePlaylistID()
	
	_, err := db.Pool.Exec(ctx, "INSERT INTO playlists (id, name, user_id) VALUES ($1, $2, $3)", id, name, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (db *Database) GetPlaylist(id string) (*Playlist, error) {
	ctx, cancel := db.ctx()
	defer cancel()

	var playlist Playlist
	err := db.Pool.QueryRow(ctx, "SELECT id, name, user_id FROM playlists WHERE id = $1", id).
		Scan(&playlist.ID, &playlist.Name, &playlist.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("playlist not found")
		}
		return nil, err
	}

	rows, err := db.Pool.Query(ctx, "SELECT url, name, track_id, duration, platform FROM playlist_songs WHERE playlist_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.URL, &song.Name, &song.TrackID, &song.Duration, &song.Platform); err == nil {
			playlist.Songs = append(playlist.Songs, song)
		}
	}
	
	return &playlist, nil
}

func (db *Database) DeletePlaylist(id string, userID int64) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, "DELETE FROM playlists WHERE id = $1 AND user_id = $2", id, userID)
	return err
}

func (db *Database) songExists(id string, trackID string) bool {
	ctx, cancel := db.ctx()
	defer cancel()

	var exists bool
	err := db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM playlist_songs WHERE playlist_id = $1 AND track_id = $2)", id, trackID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (db *Database) AddSongToPlaylist(id string, song Song) error {
	if db.songExists(id, song.TrackID) {
		return nil
	}

	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO playlist_songs (playlist_id, track_id, url, name, duration, platform) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING",
		id, song.TrackID, song.URL, song.Name, song.Duration, song.Platform)
	return err
}

func (db *Database) RemoveSongFromPlaylist(id string, trackID string) error {
	if !db.songExists(id, trackID) {
		return fmt.Errorf("track with ID %s not found in playlist", trackID)
	}

	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, "DELETE FROM playlist_songs WHERE playlist_id = $1 AND track_id = $2", id, trackID)
	return err
}

func (db *Database) GetUserPlaylists(userID int64) ([]Playlist, error) {
	ctx, cancel := db.ctx()
	defer cancel()

	rows, err := db.Pool.Query(ctx, "SELECT id, name FROM playlists WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var playlist Playlist
		playlist.UserID = userID
		if err := rows.Scan(&playlist.ID, &playlist.Name); err == nil {
			playlists = append(playlists, playlist)
		}
	}
	
	// Fetch songs for each playlist
	for i := range playlists {
		sRows, err := db.Pool.Query(ctx, "SELECT url, name, track_id, duration, platform FROM playlist_songs WHERE playlist_id = $1", playlists[i].ID)
		if err == nil {
			for sRows.Next() {
				var song Song
				if err := sRows.Scan(&song.URL, &song.Name, &song.TrackID, &song.Duration, &song.Platform); err == nil {
					playlists[i].Songs = append(playlists[i].Songs, song)
				}
			}
			sRows.Close()
		}
	}

	return playlists, nil
}

func ConvertSongsToTracks(songs []Song) []utils.MusicTrack {
	tracks := make([]utils.MusicTrack, 0, len(songs))

	for _, song := range songs {
		tracks = append(tracks, utils.MusicTrack{
			Url:      song.URL,
			Title:    song.Name,
			Id:       song.TrackID,
			Duration: song.Duration,
			Platform: song.Platform,
		})
	}

	return tracks
}
