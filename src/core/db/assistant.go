/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package db

import (
	"log/slog"
	"github.com/jackc/pgx/v5"
)

// GetAssistant retrieves the index of the assistant for a chat.
// Returns -1 if no assistant is assigned.
func (db *Database) GetAssistant(chatID int64) (int, error) {
	key := toKey(chatID)
	if cached, ok := db.assistantCache.Get(key); ok {
		return cached, nil
	}

	ctx, cancel := db.ctx()
	defer cancel()

	var num int

	err := db.Pool.QueryRow(ctx, "SELECT num FROM assistants WHERE chat_id = $1", chatID).Scan(&num)
	if err != nil {
		if err == pgx.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}

	// Cache the result
	db.assistantCache.Set(key, num)
	return num, nil
}

// SetAssistant sets the assistant index for a given chat.
func (db *Database) SetAssistant(chatID int64, num int) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx,
		"INSERT INTO assistants (chat_id, num) VALUES ($1, $2) ON CONFLICT (chat_id) DO UPDATE SET num = EXCLUDED.num",
		chatID, num)
	if err == nil {
		db.assistantCache.Set(toKey(chatID), num)
	}

	return err
}

// RemoveAssistant removes the assistant from a chat's settings.
func (db *Database) RemoveAssistant(chatID int64) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, "DELETE FROM assistants WHERE chat_id = $1", chatID)
	if err == nil {
		db.assistantCache.Delete(toKey(chatID))
	}
	return err
}

// AssignAssistant attempts to set the assistant for a chat if it is not currently set.
func (db *Database) AssignAssistant(chatID int64, proposedAssistant int) (int, error) {
	ctx, cancel := db.ctx()
	defer cancel()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	var num int
	err = tx.QueryRow(ctx, "SELECT num FROM assistants WHERE chat_id = $1 FOR UPDATE", chatID).Scan(&num)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, errx := tx.Exec(ctx, "INSERT INTO assistants (chat_id, num) VALUES ($1, $2)", chatID, proposedAssistant)
			if errx == nil {
				tx.Commit(ctx)
				db.assistantCache.Set(toKey(chatID), proposedAssistant)
				return proposedAssistant, nil
			}
			return -1, errx
		}
		return -1, err
	}

	if num == -1 {
		_, err = tx.Exec(ctx, "UPDATE assistants SET num = $1 WHERE chat_id = $2", proposedAssistant, chatID)
		if err == nil {
			tx.Commit(ctx)
			db.assistantCache.Set(toKey(chatID), proposedAssistant)
			return proposedAssistant, nil
		}
		return -1, err
	}

	tx.Commit(ctx)
	return num, nil
}

// ClearAllAssistants removes all assistant assignments.
func (db *Database) ClearAllAssistants() (int64, error) {
	ctx, cancel := db.ctx()
	defer cancel()

	tag, err := db.Pool.Exec(ctx, "DELETE FROM assistants")
	if err != nil {
		slog.Info("[DB] Error clearing assistants", "error", err)
		return 0, err
	}

	db.assistantCache.Clear()
	return tag.RowsAffected(), nil
}
