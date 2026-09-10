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
	"context"
	"log/slog"
	
	"github.com/jackc/pgx/v5"
)

type Chats struct {
	ID        int64
	PlayType  int
	AdminPlay bool
	AdminMode string
	CmdDelete bool
}

func (db *Database) getChat(chatID int64) (*Chats, error) {
	key := toKey(chatID)
	if cached, ok := db.chatCache.Get(key); ok {
		return cached, nil
	}

	var chat Chats
	ctx, cancel := db.ctx()
	defer cancel()

	err := db.Pool.QueryRow(ctx, "SELECT id, play_type, admin_play, admin_mode, cmd_delete FROM chats WHERE id = $1", chatID).
		Scan(&chat.ID, &chat.PlayType, &chat.AdminPlay, &chat.AdminMode, &chat.CmdDelete)
		
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		slog.Info("[DB] An error occurred while getting the chat", "error", err)
		return nil, err
	}

	db.chatCache.Set(key, &chat)
	return &chat, nil
}

func (db *Database) AddChat(chatID int64) error {
	chat, _ := db.getChat(chatID)
	if chat != nil {
		return nil
	}

	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, "INSERT INTO chats (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", chatID)
	if err == nil {
		slog.Info("[DB] A new chat has been added", "id", chatID)
	}
	return err
}

func (db *Database) GetPlayType(chatID int64) int {
	chat, _ := db.getChat(chatID)
	if chat == nil {
		return 0
	}
	return chat.PlayType
}

func (db *Database) SetPlayType(chatID int64, playType int) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO chats (id, play_type) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET play_type = EXCLUDED.play_type", 
		chatID, playType)
	if err == nil {
		db.chatCache.Delete(toKey(chatID))
	}
	return err
}

func (db *Database) GetPlayMode(chatID int64) bool {
	chat, _ := db.getChat(chatID)
	if chat == nil {
		return false
	}
	return chat.AdminPlay
}

func (db *Database) SetPlayMode(chatID int64, adminPlay bool) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO chats (id, admin_play) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET admin_play = EXCLUDED.admin_play", 
		chatID, adminPlay)
	if err == nil {
		db.chatCache.Delete(toKey(chatID))
	}
	return err
}

func (db *Database) GetAdminMode(chatID int64) string {
	chat, _ := db.getChat(chatID)
	if chat == nil || chat.AdminMode == "" {
		return utils.Everyone
	}
	return chat.AdminMode
}

func (db *Database) SetAdminMode(chatID int64, adminMode string) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO chats (id, admin_mode) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET admin_mode = EXCLUDED.admin_mode", 
		chatID, adminMode)
	if err == nil {
		db.chatCache.Delete(toKey(chatID))
	}
	return err
}

func (db *Database) GetCmdDelete(chatID int64) bool {
	chat, _ := db.getChat(chatID)
	if chat == nil {
		return false
	}
	return chat.CmdDelete
}

func (db *Database) SetCmdDelete(chatID int64, cmdDelete bool) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO chats (id, cmd_delete) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET cmd_delete = EXCLUDED.cmd_delete", 
		chatID, cmdDelete)
	if err == nil {
		db.chatCache.Delete(toKey(chatID))
	}
	return err
}

func (db *Database) GetAllChats() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.Pool.Query(ctx, "SELECT id, play_type, admin_play, admin_mode, cmd_delete FROM chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []int64
	for rows.Next() {
		var doc Chats
		if err := rows.Scan(&doc.ID, &doc.PlayType, &doc.AdminPlay, &doc.AdminMode, &doc.CmdDelete); err != nil {
			return nil, err
		}
		chats = append(chats, doc.ID)
		db.chatCache.Set(toKey(doc.ID), &doc)
	}
	return chats, nil
}
