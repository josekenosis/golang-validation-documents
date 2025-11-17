package database

import (
	"fmt"

	"github.com/neoway/golang-validation-documents/internal/config"
	"github.com/neoway/golang-validation-documents/internal/domain"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
		config.AppConfig.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Msg("Database connected successfully")

	if err := DB.AutoMigrate(&domain.Document{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Info().Msg("Database migrations completed")

	return nil
}

