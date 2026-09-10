/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package handlers

import (
	"simmie/src/core/cache"
	"simmie/src/vc"

	td "github.com/AshokShau/gotdbot"
)

// skipHandler handles the /skip command.
func skipHandler(c *td.Client, m *td.Message) error {
	if !adminMode(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.ReplyText(c, "The bot is not streaming in the video chat.", nil)
		return nil
	}

	_ = vc.Calls.PlayNext(c, chatID)
	return nil
}
