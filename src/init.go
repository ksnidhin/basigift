/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package src

import (
	"simmie/config"
	"simmie/src/core/db"
	"simmie/src/vc"

	"github.com/AshokShau/gotdbot"
)

func Init(client *gotdbot.Client) error {
	if err := db.InitDatabase(); err != nil {
		return err
	}

	for _, session := range config.SessionStrings {
		_, err := vc.Calls.StartClient(config.ApiId, config.ApiHash, session)
		if err != nil {
			return err
		}
	}

	vc.Calls.RegisterHandlers(client)
	return nil
}
