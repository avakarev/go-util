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
	mu             sync.Mutex
	lockingEnabled bool
	conn           *gorm.DB
	config         *gorm.Config
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

// RegisterValidationTagNameFunc registers a function to get alternate names for StructFields
func (db *DB) RegisterValidationTagNameFunc(fn validator.TagNameFunc) {
	db.validate.RegisterTagNameFunc(fn)
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

// WithLogger sets given logger as gorm logger
func (db *DB) WithLogger(l logger.Interface) *DB {
	db.config.Logger = l
	return db
}

// WithNowFunc sets given func as gorm now func
func (db *DB) WithNowFunc(fn func() time.Time) *DB {
	db.config.NowFunc = fn
	return db
}

// Open initializes db session based on dialector
func (db *DB) Open(dialector gorm.Dialector) error {
	conn, err := gorm.Open(dialector, db.config)
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

// New returns new DB value
func New() *DB {
	return &DB{validate: validator.New()}
}

// Open initializes db session based on dialector
// Deprecated: use gormutil.New().Open() instead
func Open(dialector gorm.Dialector) (*DB, error) {
	config := &gorm.Config{
		Logger: Logger,
	}
	conn, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn, config: config, validate: validator.New()}, nil
}
