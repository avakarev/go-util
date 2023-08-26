// Package gormutil implements gorm helpers
package gormutil

import (
	"gorm.io/gorm"
)

// DB defines db container
type DB struct {
	conn  *gorm.DB
	hooks *HookBus
}

// Conn returns gorm's connection
func (db *DB) Conn() *gorm.DB {
	return db.conn
}

// SubscribeHook creates subscription for create/update/delete changes of given model
func (db *DB) SubscribeHook(model interface{}, fn HookHandlerFunc) {
	if db.hooks != nil {
		db.hooks.subscribe(model, fn)
	}
}

// AfterCreateHook publishes hook after create
func (db *DB) AfterCreateHook(model interface{}) {
	if db.hooks != nil {
		db.hooks.publish(model, HookEvent(HookAfterCreate))
	}
}

// AfterUpdateHook publishes hook after update
func (db *DB) AfterUpdateHook(model interface{}) {
	if db.hooks != nil {
		db.hooks.publish(model, HookEvent(HookAfterUpdate))
	}
}

// AfterDeleteHook publishes hook after delete
func (db *DB) AfterDeleteHook(model interface{}) {
	if db.hooks != nil {
		db.hooks.publish(model, HookEvent(HookAfterDelete))
	}
}

// WithHooks enables hooks pub/sub
func (db *DB) WithHooks() *DB {
	if db.hooks == nil {
		db.hooks = newHookBus()
		go db.hooks.run()
	}
	return db
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
