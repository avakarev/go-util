// Package gormutil implements gorm helpers
package gormutil

import (
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB defines db container
type DB struct {
	mu           sync.Mutex
	locksEnabled bool
	conn         *gorm.DB
	config       *gorm.Config
	validate     *validator.Validate
	hooks        *HookBus
}

// ConfigureFunc defines configurator func
type ConfigureFunc func(*DB) error

// Conn returns gorm's connection
func (db *DB) Conn() *gorm.DB {
	return db.conn
}

// Begin begins a transaction
func (db *DB) Begin() *DB {
	return &DB{
		locksEnabled: db.locksEnabled,
		conn:         db.conn.Begin(),
		config:       db.config,
		validate:     db.validate,
		hooks:        db.hooks,
	}
}

// Rollback rollbacks the transaction
func (db *DB) Rollback() {
	db.conn.Rollback()
}

// Commit commits the transaction
func (db *DB) Commit() error {
	return db.conn.Commit().Error
}

// RegisterValidation adds a custom validation for the given tag
func (db *DB) RegisterValidation(tag string, fn validator.Func) error {
	return db.validate.RegisterValidation(tag, fn)
}

// RegisterValidationTagNameFunc registers a function to get alternate names for StructFields
func (db *DB) RegisterValidationTagNameFunc(fn validator.TagNameFunc) {
	db.validate.RegisterTagNameFunc(fn)
}

// SubscribeHook creates subscription for create/update/delete changes of given model
func (db *DB) SubscribeHook(model any, fn HookHandlerFunc) {
	if db.hooks != nil {
		db.hooks.subscribe(model, fn)
	}
}

// AfterCreateHook publishes hook after create
func (db *DB) AfterCreateHook(model any) {
	if db.hooks != nil {
		db.hooks.publish(model, HookEvent(HookAfterCreate))
	}
}

// AfterUpdateHook publishes hook after update
func (db *DB) AfterUpdateHook(model any) {
	if db.hooks != nil {
		db.hooks.publish(model, HookEvent(HookAfterUpdate))
	}
}

// AfterDeleteHook publishes hook after delete
func (db *DB) AfterDeleteHook(model any) {
	if db.hooks != nil {
		db.hooks.publish(model, HookEvent(HookAfterDelete))
	}
}

// WithHooks enables hooks pub/sub
func (db *DB) WithHooks() {
	if db.hooks == nil {
		db.hooks = newHookBus()
		go db.hooks.run()
	}
}

// WithLocks enables mutex locks during create/update/delete calls
func WithLocks() ConfigureFunc {
	return func(db *DB) error {
		db.locksEnabled = true
		return nil
	}
}

// WithLogger sets given logger as gorm logger
func WithLogger(l logger.Interface) ConfigureFunc {
	return func(db *DB) error {
		db.config.Logger = l
		return nil
	}
}

// WithNowFunc sets given func as gorm now func
func WithNowFunc(fn func() time.Time) ConfigureFunc {
	return func(db *DB) error {
		db.config.NowFunc = fn
		return nil
	}
}

// Open initializes db session based on dialector
func Open(dialector gorm.Dialector, fns ...ConfigureFunc) (*DB, error) {
	db := &DB{config: &gorm.Config{}, validate: validator.New()}
	for _, fn := range fns {
		if err := fn(db); err != nil {
			return nil, err
		}
	}
	conn, err := gorm.Open(dialector, db.config)
	if err != nil {
		return nil, err
	}
	db.conn = conn
	return db, nil
}
