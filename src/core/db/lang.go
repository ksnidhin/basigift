/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (db *Database) GetLanguage(chatID int64) (string, error) {
	key := toKey(chatID)
	if cached, ok := db.langCache.Get(key); ok {
		return cached, nil
	}

	ctx, cancel := db.ctx()
	defer cancel()

	var lang string
	err := db.Pool.QueryRow(ctx, "SELECT lang_code FROM chat_languages WHERE chat_id = $1", chatID).Scan(&lang)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "en", nil
		}
		return "", err
	}
	db.langCache.Set(key, lang)
	return lang, nil
}

func (db *Database) SetLanguage(ctx context.Context, chatID int64, langCode string) error {
	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO chat_languages (chat_id, lang_code) VALUES ($1, $2) ON CONFLICT (chat_id) DO UPDATE SET lang_code = EXCLUDED.lang_code", 
		chatID, langCode)

	if err == nil {
		db.langCache.Set(toKey(chatID), langCode)
	}
	return err
}
