package database

import (
	"context"
	"time"
	"usermanagement-api/config"
	"usermanagement-api/domain/entities"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// ConnectDB connects to the PostgreSQL database
func ConnectDB(cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()
	zapLogger.Debug("Connecting to database", zap.String("host", cfg.Database.Host), zap.Int("port", cfg.Database.Port))

	// Map config log level to GORM log level
	gormLogLevel := gormLogger.Info // default
	switch cfg.Database.LogLevel {
	case "silent":
		gormLogLevel = gormLogger.Silent
	case "error":
		gormLogLevel = gormLogger.Error
	case "warn":
		gormLogLevel = gormLogger.Warn
	case "info":
		gormLogLevel = gormLogger.Info
	}

	// Create GORM logger adapter for zap
	gormZapLogger := &gormZapLogger{
		logger: zapLogger.With(zap.String("component", "gorm")),
		level:  gormLogLevel,
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormZapLogger,
	})
	if err != nil {
		zapLogger.Error("Failed to connect to database", zap.Error(err), zap.String("dsn", maskDSN(dsn)))
		return nil, err
	}

	zapLogger.Info("Connected to database", zap.String("host", cfg.Database.Host), zap.Int("port", cfg.Database.Port))
	return db, nil
}

// MigrateDB migrates the database schema
func MigrateDB(db *gorm.DB, zapLogger *zap.Logger) error {
	zapLogger.Info("Starting database migration")

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		zapLogger.Warn("Failed to create uuid extension (might already exist)", zap.Error(err))
	}

	err := db.AutoMigrate(
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
		&entities.Menu{},
		&entities.ModelPermission{},
		&entities.Setting{},
		&entities.UserMeta{},
	)
	if err != nil {
		zapLogger.Error("Failed to migrate database", zap.Error(err))
		return err
	}

	zapLogger.Info("Database migration completed")
	return nil
}

// gormZapLogger implements gorm.io/gorm/logger.Interface
type gormZapLogger struct {
	logger *zap.Logger
	level  gormLogger.LogLevel
}

func (l *gormZapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *gormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Info {
		l.logger.Sugar().Infof(msg, data...)
	}
}

func (l *gormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Warn {
		l.logger.Sugar().Warnf(msg, data...)
	}
}

func (l *gormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Error {
		l.logger.Sugar().Errorf(msg, data...)
	}
}

func (l *gormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.level >= gormLogger.Error:
		// Log errors at Error level
		l.logger.Error("SQL query error",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	case elapsed > 200*time.Millisecond && l.level >= gormLogger.Warn:
		// Log slow queries at Warn level
		l.logger.Warn("Slow SQL query",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	case l.level >= gormLogger.Info:
		// Log all queries at Info level when GORM log level is Info or lower
		// Use Info level instead of Debug so it shows even when zap logger is at Info level
		l.logger.Info("SQL query",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	}
}

// maskDSN masks password in DSN string for logging
func maskDSN(dsn string) string {
	// Simple masking - in production you might want more sophisticated masking
	if len(dsn) > 50 {
		return dsn[:30] + "***masked***"
	}
	return "***masked***"
}
