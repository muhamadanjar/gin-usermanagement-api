package database

import (
	"fmt"
	"log"
	"os"
	"usermanagement-api/domain/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB connects to the PostgreSQL database
func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")
	return db
}

// MigrateDB migrates the database schema
func MigrateDB(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err := db.AutoMigrate(
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
		&entities.Menu{},
		&entities.ModelPermission{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
}
