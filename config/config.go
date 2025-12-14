package config

import (
	"fmt"
	"strings"
	"usermanagement-api/pkg/logger"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	Redis    RedisConfig
	Firebase FirebaseConfig
	CORS     CORSConfig
	Logger   logger.Config
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	LogLevel string // silent, error, warn, info
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret                string
	RefreshTokenSecret    string
	Expiration            int // in hours
	AccessTokenExpiration int // in hours
	RefreshTokenExpiration int // in hours
}

// EmailConfig holds email-related configuration
type EmailConfig struct {
	Host         string
	Port         int
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// FirebaseConfig holds Firebase-related configuration
type FirebaseConfig struct {
	CredentialsFile string
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	MaxAge           int // in seconds
}

// LoadConfig loads configuration using viper
// Priority: .env file > environment variables > config files (yaml, json, toml)
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set config file name (without extension)
	v.SetConfigName(".env")
	v.SetConfigType("env")

	// Add paths to search for config files
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Enable environment variable reading
	v.AutomaticEnv()

	// Replace dots with underscores for env vars (e.g., database.host -> DATABASE_HOST)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read .env file (optional - will continue if not found)
	if err := v.ReadInConfig(); err != nil {
		// If .env file doesn't exist, that's okay - will use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Also try to read other config formats (yaml, json, toml) if they exist
	// These will override .env if they exist
	for _, ext := range []string{"yaml", "yml", "json", "toml"} {
		v.SetConfigType(ext)
		if err := v.MergeInConfig(); err == nil {
			// Successfully loaded additional config file
			break
		}
	}

	// Load server config
	port := v.GetInt("server.port")
	if port == 0 {
		port = v.GetInt("SERVER_PORT")
		if port == 0 {
			port = 8080
		}
	}

	host := v.GetString("server.host")
	if host == "" {
		host = v.GetString("SERVER_HOST")
		if host == "" {
			host = "0.0.0.0"
		}
	}

	// Load database config
	dbPort := v.GetInt("database.port")
	if dbPort == 0 {
		dbPort = v.GetInt("DB_PORT")
		if dbPort == 0 {
			dbPort = 5432
		}
	}

	dbHost := v.GetString("database.host")
	if dbHost == "" {
		dbHost = v.GetString("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}
	}

	dbUser := v.GetString("database.user")
	if dbUser == "" {
		dbUser = v.GetString("DB_USER")
		if dbUser == "" {
			dbUser = "postgres"
		}
	}

	dbPassword := v.GetString("database.password")
	if dbPassword == "" {
		dbPassword = v.GetString("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = "postgres"
		}
	}

	dbName := v.GetString("database.name")
	if dbName == "" {
		dbName = v.GetString("DB_NAME")
		if dbName == "" {
			dbName = "user_management"
		}
	}

	dbSSLMode := v.GetString("database.sslmode")
	if dbSSLMode == "" {
		dbSSLMode = v.GetString("DB_SSLMODE")
		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}
	}

	// Load database log level (for GORM SQL query logging)
	dbLogLevel := v.GetString("database.log_level")
	if dbLogLevel == "" {
		dbLogLevel = v.GetString("DB_LOG_LEVEL")
		if dbLogLevel == "" {
			dbLogLevel = "info" // default: only show errors and slow queries
		}
	}

	// Load JWT config
	jwtSecret := v.GetString("jwt.secret")
	if jwtSecret == "" {
		jwtSecret = v.GetString("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your_jwt_secret_key"
		}
	}

	refreshTokenSecret := v.GetString("jwt.refresh_token_secret")
	if refreshTokenSecret == "" {
		refreshTokenSecret = v.GetString("REFRESH_TOKEN_SECRET")
		if refreshTokenSecret == "" {
			refreshTokenSecret = "your_refresh_token_secret_key"
		}
	}

	jwtExpiration := v.GetInt("jwt.expiration")
	if jwtExpiration == 0 {
		jwtExpiration = v.GetInt("JWT_EXPIRATION")
		if jwtExpiration == 0 {
			jwtExpiration = 24
		}
	}

	accessTokenExpiration := v.GetInt("jwt.access_token_expiration")
	if accessTokenExpiration == 0 {
		accessTokenExpiration = v.GetInt("ACCESS_TOKEN_EXPIRATION")
		if accessTokenExpiration == 0 {
			accessTokenExpiration = 1 // 1 hour default
		}
	}

	refreshTokenExpiration := v.GetInt("jwt.refresh_token_expiration")
	if refreshTokenExpiration == 0 {
		refreshTokenExpiration = v.GetInt("REFRESH_TOKEN_EXPIRATION")
		if refreshTokenExpiration == 0 {
			refreshTokenExpiration = 168 // 7 days default
		}
	}

	// Load Redis config
	redisAddr := v.GetString("redis.addr")
	if redisAddr == "" {
		redisAddr = v.GetString("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}
	}

	redisPassword := v.GetString("redis.password")
	if redisPassword == "" {
		redisPassword = v.GetString("REDIS_PASSWORD")
	}

	redisDB := v.GetInt("redis.db")
	if redisDB == 0 {
		redisDB = v.GetInt("REDIS_DB")
	}

	// Load Firebase config
	firebaseCredentialsFile := v.GetString("firebase.credentials_file")
	if firebaseCredentialsFile == "" {
		firebaseCredentialsFile = v.GetString("FIREBASE_CREDENTIALS_FILE")
	}

	// Load CORS config
	corsOriginsStr := v.GetString("cors.allowed_origins")
	if corsOriginsStr == "" {
		corsOriginsStr = v.GetString("CORS_ALLOWED_ORIGINS")
	}
	corsOrigins := []string{"http://localhost:3000"} // default
	if corsOriginsStr != "" {
		corsOrigins = strings.Split(corsOriginsStr, ",")
		for i, origin := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(origin)
		}
	}

	allowCredentials := v.GetBool("cors.allow_credentials")
	if !v.IsSet("cors.allow_credentials") {
		allowCredentialsStr := v.GetString("CORS_ALLOW_CREDENTIALS")
		if allowCredentialsStr == "false" {
			allowCredentials = false
		} else {
			allowCredentials = true // default
		}
	}

	corsMaxAge := v.GetInt("cors.max_age")
	if corsMaxAge == 0 {
		corsMaxAge = v.GetInt("CORS_MAX_AGE")
		if corsMaxAge == 0 {
			corsMaxAge = 43200 // 12 hours default
		}
	}

	// Load Email config
	emailHost := v.GetString("email.host")
	if emailHost == "" {
		emailHost = v.GetString("SMTP_HOST")
		if emailHost == "" {
			emailHost = "smtp.example.com"
		}
	}

	emailPort := v.GetInt("email.port")
	if emailPort == 0 {
		emailPort = v.GetInt("SMTP_PORT")
		if emailPort == 0 {
			emailPort = 587
		}
	}

	emailSenderName := v.GetString("email.sender_name")
	if emailSenderName == "" {
		emailSenderName = v.GetString("SMTP_SENDER_NAME")
		if emailSenderName == "" {
			emailSenderName = "Your App"
		}
	}

	emailAuthEmail := v.GetString("email.auth_email")
	if emailAuthEmail == "" {
		emailAuthEmail = v.GetString("SMTP_AUTH_EMAIL")
	}

	emailAuthPassword := v.GetString("email.auth_password")
	if emailAuthPassword == "" {
		emailAuthPassword = v.GetString("SMTP_AUTH_PASSWORD")
	}

	// Load Logger config
	loggerLevel := v.GetString("logger.level")
	if loggerLevel == "" {
		loggerLevel = v.GetString("LOGGER_LEVEL")
		if loggerLevel == "" {
			loggerLevel = "info"
		}
	}

	loggerMode := v.GetString("logger.mode")
	if loggerMode == "" {
		loggerMode = v.GetString("LOGGER_MODE")
		if loggerMode == "" {
			loggerMode = "development"
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Host: host,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
			LogLevel: dbLogLevel,
		},
		JWT: JWTConfig{
			Secret:                jwtSecret,
			RefreshTokenSecret:    refreshTokenSecret,
			Expiration:            jwtExpiration,
			AccessTokenExpiration: accessTokenExpiration,
			RefreshTokenExpiration: refreshTokenExpiration,
		},
		Email: EmailConfig{
			Host:         emailHost,
			Port:         emailPort,
			SenderName:   emailSenderName,
			AuthEmail:    emailAuthEmail,
			AuthPassword: emailAuthPassword,
		},
		Redis: RedisConfig{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
		Firebase: FirebaseConfig{
			CredentialsFile: firebaseCredentialsFile,
		},
		CORS: CORSConfig{
			AllowedOrigins:   corsOrigins,
			AllowCredentials: allowCredentials,
			MaxAge:           corsMaxAge,
		},
		Logger: logger.Config{
			Level: loggerLevel,
			Mode:  loggerMode,
		},
	}, nil
}

// GetDSN returns the database DSN (Data Source Name)
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}