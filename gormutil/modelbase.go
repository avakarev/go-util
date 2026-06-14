package gormutil

import (
	"time"

	"github.com/google/uuid"
)

// ModelBase contains common columns for all tables
type ModelBase struct {
	ID        uuid.UUID `json:"id" yaml:"id" validate:"required" gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// GenerateID generates and assigns new id value
func (base *ModelBase) GenerateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	base.ID = id
	return nil
}
