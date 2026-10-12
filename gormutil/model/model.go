// Package model provides embeddable base structs for gorm models.
//
// Primary keys are assigned via the exported EnsureID method rather than a
// gorm BeforeCreate hook. A hook declared on an embedded base is promoted to
// the outer model, so an end model that defines its own BeforeCreate would
// silently shadow the base one and skip id generation. Keeping EnsureID as a
// plain method avoids that trap: gormutil.Create invokes it through the Base
// interface, independent of the model's method set, and end models remain free
// to define their own gorm hooks without interfering with id assignment.
package model

import (
	"time"

	"github.com/google/uuid"
)

// Base is implemented by models that can lazily assign their primary key.
// Default and Timeseries satisfy it.
type Base interface {
	EnsureID() error
}

// Default contains common columns for all tables: a uuid primary key plus
// created_at / updated_at timestamps.
type Default struct {
	ID        uuid.UUID `json:"id" yaml:"id" validate:"required" gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// EnsureID assigns a newly generated id when the current one is nil. It's a
// no-op when the id is already set, so callers may provide their own id.
func (m *Default) EnsureID() error {
	if m.ID != uuid.Nil {
		return nil
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

// Timeseries contains common columns for time series tables: a uuid primary key
// plus an indexed created_at timestamp. Unlike Default it has no updated_at
// column, since time series rows are append-only.
type Timeseries struct {
	ID        uuid.UUID `json:"id" yaml:"id" validate:"required" gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt" gorm:"index"`
}

// EnsureID assigns a newly generated id when the current one is nil. It's a
// no-op when the id is already set, so callers may provide their own id.
func (m *Timeseries) EnsureID() error {
	if m.ID != uuid.Nil {
		return nil
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}
