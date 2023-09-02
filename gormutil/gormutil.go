// Package gormutil implements gorm helpers
package gormutil

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// DB defines db container
type DB struct {
	mu             sync.Mutex
	lockingEnabled bool
	conn           *gorm.DB
	validate       *validator.Validate
	hooks          *HookBus
}

// Conn returns gorm's connection
func (db *DB) Conn() *gorm.DB {
	return db.conn
}

// RegisterValidation adds a custom validation for the given tag
func (db *DB) RegisterValidation(tag string, fn validator.Func) error {
	return db.validate.RegisterValidation(tag, fn)
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

// WithLocking enables mutex locking during create/update/delete calls
func (db *DB) WithLocking() *DB {
	db.lockingEnabled = true
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
	return &DB{conn: conn, validate: validator.New()}, nil
}
