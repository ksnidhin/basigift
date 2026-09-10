/*
 * Daddy Noah - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/Simmie/DaddyNoah
 */

package db

// AddBlacklistedChat adds a chat to the blacklist.
func (db *Database) AddBlacklistedChat(chatID int64) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx,
		"INSERT INTO blacklisted_chats (chat_id) VALUES ($1) ON CONFLICT (chat_id) DO NOTHING",
		chatID)
	if err == nil {
		db.blChatsCache.Delete("bl_chats")
	}
	return err
}

// RemoveBlacklistedChat removes a chat from the blacklist.
func (db *Database) RemoveBlacklistedChat(chatID int64) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, "DELETE FROM blacklisted_chats WHERE chat_id = $1", chatID)
	if err == nil {
		db.blChatsCache.Delete("bl_chats")
	}
	return err
}

// GetBlacklistedChats retrieves the list of blacklisted chat IDs.
func (db *Database) GetBlacklistedChats() []int64 {
	if cached, ok := db.blChatsCache.Get("bl_chats"); ok {
		return cached
	}
	
	ctx, cancel := db.ctx()
	defer cancel()

	rows, err := db.Pool.Query(ctx, "SELECT chat_id FROM blacklisted_chats")
	if err != nil {
		return []int64{}
	}
	defer rows.Close()

	var chats []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			chats = append(chats, id)
		}
	}
	
	db.blChatsCache.Set("bl_chats", chats)
	return chats
}

// IsBlacklistedChat checks if a chat is blacklisted.
func (db *Database) IsBlacklistedChat(chatID int64) bool {
	chats := db.GetBlacklistedChats()
	return contains(chats, chatID)
}

// AddBlacklistedUser adds a user to the blacklist.
func (db *Database) AddBlacklistedUser(userID int64) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx,
		"INSERT INTO blacklisted_users (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING",
		userID)
	if err == nil {
		db.blUsersCache.Delete("bl_users")
	}
	return err
}

// RemoveBlacklistedUser removes a user from the blacklist.
func (db *Database) RemoveBlacklistedUser(userID int64) error {
	ctx, cancel := db.ctx()
	defer cancel()

	_, err := db.Pool.Exec(ctx, "DELETE FROM blacklisted_users WHERE user_id = $1", userID)
	if err == nil {
		db.blUsersCache.Delete("bl_users")
	}
	return err
}

// GetBlacklistedUsers retrieves the list of blacklisted user IDs.
func (db *Database) GetBlacklistedUsers() []int64 {
	if cached, ok := db.blUsersCache.Get("bl_users"); ok {
		return cached
	}
	
	ctx, cancel := db.ctx()
	defer cancel()

	rows, err := db.Pool.Query(ctx, "SELECT user_id FROM blacklisted_users")
	if err != nil {
		return []int64{}
	}
	defer rows.Close()

	var users []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			users = append(users, id)
		}
	}
	
	db.blUsersCache.Set("bl_users", users)
	return users
}

// IsBlacklistedUser checks if a user is blacklisted.
func (db *Database) IsBlacklistedUser(userID int64) bool {
	users := db.GetBlacklistedUsers()
	return contains(users, userID)
}
