// Package gormutil implements gorm helpers
package gormutil

import (
	"gorm.io/gorm"
)

// DB defines db container
type DB struct {
	conn *gorm.DB
}

// Conn returns gorm's connection
func (db *DB) Conn() *gorm.DB {
	return db.conn
}

// Open initialize db session based on dialector
func Open(dialector gorm.Dialector) (*DB, error) {
	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger: Logger,
	})
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}
