package db

import (
	"usermanagement-api/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Menu{}, &models.ModelPermission{})
}
