package gormutil

import (
	"time"

	"github.com/google/uuid"
)

// ModelBase contains common columns for all tables
//
// Deprecated: use model.Default instead. ModelBase is kept for backward
// compatibility and will be removed in a future release.
type ModelBase struct {
	ID        uuid.UUID `json:"id" yaml:"id" validate:"required" gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// GenerateID generates and assigns new id value
//
// Deprecated: use model.Default and its EnsureID method instead. GenerateID is
// kept for backward compatibility and will be removed in a future release.
func (base *ModelBase) GenerateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	base.ID = id
	return nil
}
