package gormutil

import (
	"time"

	"github.com/google/uuid"
)

// ModelBase contains common columns for all tables
type ModelBase struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex" json:"id" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
