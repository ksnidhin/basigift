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
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"simmie/config"
	"simmie/src/core/cache"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool

	chatCache      *cache.Cache[*Chats]
	userCache      *cache.Cache[
Users]
	assistantCache *cache.Cache[int]
	authCache      *cache.Cache[[]int64]
	langCache      *cache.Cache[string]
	loggerCache    *cache.Cache[bool]
	blChatsCache   *cache.Cache[[]int64]
	blUsersCache   *cache.Cache[[]int64]
}

var Instance *Database

func InitDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgxConfig, err := pgxpool.ParseConfig(config.DatabaseUrl)
	if err != nil {
		return err
	}

	pgxConfig.MaxConns = config.DbMaxConns
	pgxConfig.MinConns = config.DbMinConns

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return errors.New("failed to ping postgres database: " + err.Error())
	}

	Instance = &Database{
		Pool:           pool,
		chatCache:      cache.NewCache[*Chats](60 * time.Minute),
		userCache:      cache.NewCache[*Users](60 * time.Minute),
		assistantCache: cache.NewCache[int](2 * time.Hour),
		authCache:      cache.NewCache[[]int64](20 * time.Minute),
		langCache:      cache.NewCache[string](20 * time.Minute),
		loggerCache:    cache.NewCache[bool](24 * time.Hour),
		blChatsCache:   cache.NewCache[[]int64](20 * time.Minute),
		blUsersCache:   cache.NewCache[[]int64](20 * time.Minute),
	}

	if err := Instance.RunMigrations(); err != nil {
		return errors.New("failed to run database migrations: " + err.Error())
	}

	slog.Info("[DB] PostgreSQL database connection established.")
	return nil
}

func (db *Database) RunMigrations() error {
	ctx, cancel := db.ctx()
	defer cancel()

	sqlPath := filepath.Join("src", "core", "db", "migrations", "001_initial.sql")
	hederContent, err := os.ReadFile(sqlPath)
	if err != nil {
		sqlPath = filepath.Join("migrations", "001_initial.sql")
		hederContent, err = os.ReadFile(sqlPath)
		if err != nil {
			slog.Warn("[DB] Migrations file not found, skipping.", "err", err)
			return nil
		}
	}

	queries := strings.Split(string(hederContent), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := db.Pool.Exec(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Close() error {
	slog.Info("[DB] Closing PostgreSQL connection pool...")
	db.Pool.Close()
	return nil
}

func (db *Database) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (db *Database) Ping() error {
	ctx, cancel := db.ctx()
	defer cancel()
	return db.Pool.Ping(ctx)
}
