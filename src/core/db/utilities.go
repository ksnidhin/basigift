/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package db

import (
	"fmt"
)

// toKey converts an int64 to a string key.
func toKey(id int64) string {
	return fmt.Sprintf("%d", id)
}

// contains checks if a slice contains a specific int64 value.
func contains(slice []int64, item int64) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
