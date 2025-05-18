package database

import (
	"github.com/golgoth31/aliasme/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Config holds the database configuration
type Config struct {
	Path string
}

// New creates a new database connection
func New(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Email{}, &models.Alias{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to migrate database")
		return nil, err
	}

	return db, nil
}
