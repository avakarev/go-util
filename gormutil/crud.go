package gormutil

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// First returns first row matching the given query
func First[T any](tx *gorm.DB) *T {
	var model T
	if err := tx.First(&model).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) && tx.Logger != nil {
			tx.Logger.Error(context.Background(), "failed to query database, got error %v", err)
		}
		return nil
	}
	return &model
}

// Find returns all rows models by the given query
func Find[T any](tx *gorm.DB) []T {
	var models []T
	if err := tx.Find(&models).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) && tx.Logger != nil {
			tx.Logger.Error(context.Background(), "failed to query database, got error %v", err)
		}
		return nil
	}
	return models
}

// Count returns number of record in given table
func (db *DB) Count(model any) int64 {
	var count int64
	db.Conn().Model(model).Count(&count)
	return count
}

// CountBy returns number of record in given table with given conditions
func (db *DB) CountBy(model any, cond any, args ...any) int64 {
	var count int64
	db.Conn().Model(model).Where(cond, args).Count(&count)
	return count
}

// ExistsBy checks whether given model exists with given conditions
//
// @TODO: try to optimize the query to something like
// SELECT EXISTS(SELECT 1 FROM vaults WHERE id="foobar" LIMIT 1);
func (db *DB) ExistsBy(model any, cond any, args ...any) bool {
	var exists bool
	q := db.Conn().
		Model(model).Select("count(*) > 0").
		Where(cond, args)
	if err := q.Find(&exists).Error; err != nil && db.config.Logger != nil {
		db.config.Logger.Error(context.Background(), "failed to query database, got error %v", err)
	}
	return exists
}

// ExistsByID checks whether given model exists with given id
func (db *DB) ExistsByID(model any, id string) bool {
	return db.ExistsBy(model, "id = ?", id)
}

// Validate validates given model struct
func (db *DB) Validate(model any) error {
	return db.validate.Struct(model)
}

// Create validates and persists new record
func (db *DB) Create(model any) error {
	if db.locksEnabled {
		db.mu.Lock()
		defer db.mu.Unlock()
	}

	if err := db.Validate(model); err != nil {
		return err
	}
	if err := db.Conn().Create(model).Error; err != nil {
		return err
	}

	db.AfterCreateHook(model)
	return nil
}

// Changeset extracts values of given field names from the model
func Changeset(model any, names []string) (map[string]any, error) {
	data := make(map[string]any)
	source := reflect.ValueOf(model)
	if source.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("model is expected to be <ptr>, instead <%T> is given", model)
	}
	source = source.Elem() // dereference the ptr
	if source.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model is expected to be <struct>, instead <%s> is given", source.Kind())
	}
	ns := schema.NamingStrategy{}
	for _, n := range names {
		f := source.FieldByName(n)
		if !f.IsValid() {
			return nil, fmt.Errorf("model doesn't have %s field", n)
		}
		data[ns.ColumnName("", n)] = f.Interface()
	}
	return data, nil
}

// Update validates and persists existing record
func (db *DB) Update(model any, names ...string) error {
	if db.locksEnabled {
		db.mu.Lock()
		defer db.mu.Unlock()
	}

	if err := db.Validate(model); err != nil {
		return err
	}

	if len(names) == 0 {
		return db.Conn().Updates(model).Error
	}

	data, err := Changeset(model, names)
	if err != nil {
		return err
	}

	if err := db.Conn().Model(model).Updates(data).Error; err != nil {
		return err
	}

	db.AfterUpdateHook(model)
	return nil
}

// Delete deletes given record from the db table
func (db *DB) Delete(model any, conds ...any) error {
	if db.locksEnabled {
		db.mu.Lock()
		defer db.mu.Unlock()
	}

	if err := db.Conn().Delete(model, conds...).Error; err != nil {
		return err
	}

	db.AfterDeleteHook(model)
	return nil
}

// DeleteByID deletes given record with given id from the db table
func (db *DB) DeleteByID(model any, id string) error {
	source := reflect.ValueOf(model)
	if source.Kind() != reflect.Pointer {
		return fmt.Errorf("model is expected to be <ptr>, instead <%T> is given", model)
	}
	source = source.Elem() // dereference the ptr
	if source.Kind() != reflect.Struct {
		return fmt.Errorf("model is expected to be <struct>, instead <%s> is given", source.Kind())
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	idField := source.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		idField.Set(reflect.ValueOf(uid))
	} else {
		return fmt.Errorf("can't set ID field for <%T> model", model)
	}

	return db.Delete(model)
}
