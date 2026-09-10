/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package db

import (
	"github.com/jackc/pgx/v5"
)

func (db *Database) GetLoggerStatus() bool {
	if cached, ok := db.loggerCache.Get("logger"); ok {
		return cached
	}

	ctx, cancel := db.ctx()
	defer cancel()

	var val string
	err := db.Pool.QueryRow(ctx, "SELECT value FROM bot_settings WHERE key = 'logger'").Scan(&val)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false
		}
		return false
	}
	
	status := (val == "true")
	db.loggerCache.Set("logger", status)
	return status
}

func (db *Database) SetLoggerStatus(status bool) error {
	ctx, cancel := db.ctx()
	defer cancel()
	
	val := "false"
	if status {
		val = "true"
	}
	
	_, err := db.Pool.Exec(ctx, 
		"INSERT INTO bot_settings (key, value) VALUES ('logger', $1) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value", 
		val)
	if err == nil {
		db.loggerCache.Set("logger", status)
	}
	return err
}
