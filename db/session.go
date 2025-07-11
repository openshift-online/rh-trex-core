package db

import (
	"context"
	"gorm.io/gorm"
)

// SessionFactory defines the interface for creating database sessions
type SessionFactory interface {
	New(ctx context.Context) *gorm.DB
}

// BasicSessionFactory provides a simple session factory implementation
type BasicSessionFactory struct {
	db *gorm.DB
}

// NewBasicSessionFactory creates a new basic session factory
func NewBasicSessionFactory(db *gorm.DB) *BasicSessionFactory {
	return &BasicSessionFactory{db: db}
}

// New creates a new database session
func (f *BasicSessionFactory) New(ctx context.Context) *gorm.DB {
	return f.db.WithContext(ctx)
}