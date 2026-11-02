package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to SQLite and brings the schema up to date. SQLite takes one
// writer, so the pool is a single connection and the database cannot be split
// across two of them.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", dsn, err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("connection pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&Member{}); err != nil {
		return nil, fmt.Errorf("migrating: %w", err)
	}
	return db, nil
}
