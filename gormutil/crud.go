package gormutil

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var validate *validator.Validate

// First returns first row matching the given query
func First[T any](tx *gorm.DB) *T {
	var model T
	if err := tx.First(&model).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("error while querying db")
		}
		return nil
	}
	return &model
}

// Find returns all rows models by the given query
func Find[T any](tx *gorm.DB) []T {
	var models []T
	if err := tx.Find(&models).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("error while querying db")
		}
		return nil
	}
	return models
}

// Count returns number of record in given table
func (db *DB) Count(model interface{}) int64 {
	var count int64
	db.Conn().Model(model).Count(&count)
	return count
}

// CountBy returns number of record in given table with given conditions
func (db *DB) CountBy(model interface{}, cond interface{}, args ...interface{}) int64 {
	var count int64
	db.Conn().Model(model).Where(cond, args).Count(&count)
	return count
}

// ExistsBy checks whether given model exists with given conditions
//
// @TODO: try to optimize the query to something like
// SELECT EXISTS(SELECT 1 FROM vaults WHERE id="foobar" LIMIT 1);
func (db *DB) ExistsBy(model interface{}, cond interface{}, args ...interface{}) bool {
	var exists bool
	q := db.Conn().
		Model(model).Select("count(*) > 0").
		Where(cond, args)
	if err := q.Find(&exists).Error; err != nil {
		log.Error().Err(err).Send()
	}
	return exists
}

// ExistsByID checks whether given model exists with given id
func (db *DB) ExistsByID(model interface{}, id string) bool {
	return db.ExistsBy(model, "id = ?", id)
}

// Create validates and persists new record
func (db *DB) Create(model interface{}) error {
	if err := validate.Struct(model); err != nil {
		return err
	}
	return db.Conn().Create(model).Error
}

// Changeset extracts values of given field names from the model
func Changeset(model interface{}, names []string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	source := reflect.ValueOf(model)
	if source.Kind() != reflect.Ptr {
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
func (db *DB) Update(model interface{}, names ...string) error {
	if err := validate.Struct(model); err != nil {
		return err
	}

	if len(names) == 0 {
		return db.Conn().Updates(model).Error
	}

	data, err := Changeset(model, names)
	if err != nil {
		return err
	}
	return db.Conn().Model(model).Updates(data).Error
}

func init() {
	validate = validator.New()
}
