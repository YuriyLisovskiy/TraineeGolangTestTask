package app

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetenvOrDefault(key, default_ string) string {
	val := os.Getenv(key)
	if val == "" {
		return default_
	}

	return val
}

func ConnectToPostgreSQLWithEnv() (*gorm.DB, error) {
	port, err := strconv.Atoi(GetenvOrDefault("POSTGRES_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		GetenvOrDefault("POSTGRES_HOST", "localhost"),
		GetenvOrDefault("POSTGRES_USER", "postgres"),
		GetenvOrDefault("POSTGRES_PASSWORD", ""),
		GetenvOrDefault("POSTGRES_DB_NAME", ""),
		port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	return db, nil
}
