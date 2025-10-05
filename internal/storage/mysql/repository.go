package mysql

import (
	"context"

	"gorm.io/gorm"
)

// Repository exposes helpers for CRUD operations on the MySQL storage.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a repository using the provided GORM connection.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DB returns the underlying GORM instance.
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// Create inserts the provided entity.
func (r *Repository) Create(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update persists changes of the provided entity.
func (r *Repository) Update(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// FindByID loads a single entity by ULID.
func (r *Repository) FindByID(ctx context.Context, dest interface{}, id string) error {
	return r.db.WithContext(ctx).First(dest, "id = ?", id).Error
}

// List returns all records for the given destination slice pointer.
func (r *Repository) List(ctx context.Context, dest interface{}) error {
	return r.db.WithContext(ctx).Find(dest).Error
}

// DeleteByID removes an entity by its ULID.
func (r *Repository) DeleteByID(ctx context.Context, model interface{}, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(model).Error
}
