package database

import (
	"usermanagement-api/infrastructure/database/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateDB migrates the database schema. It lives here (infrastructure), not
// in pkg/database, because it knows the concrete GORM models.
func MigrateDB(db *gorm.DB, zapLogger *zap.Logger) error {
	zapLogger.Info("Starting database migration")

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		zapLogger.Warn("Failed to create uuid extension (might already exist)", zap.Error(err))
	}

	err := db.AutoMigrate(
		&models.UserModel{},
		&models.RoleModel{},
		&models.PermissionModel{},
		&models.MenuModel{},
		&models.MenuPermissionLink{},
		&models.SettingModel{},
		&models.UserMetaModel{},
		&models.TokenHistoryModel{},
	)
	if err != nil {
		zapLogger.Error("Failed to migrate database", zap.Error(err))
		return err
	}

	zapLogger.Info("Database migration completed")
	return nil
}
