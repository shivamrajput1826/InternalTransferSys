package database

import (
	"fmt"
	"internal-transfer-system/config"
	"internal-transfer-system/internal/models"
	customLogger "internal-transfer-system/logger"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	DB       *gorm.DB
	dbLogger = customLogger.CreateLogger("database")
	once     sync.Once
)

func InitDB() error {
	var initErr error
	once.Do(func() {
		initErr = initializeDB()
	})
	return initErr
}

func initializeDB() error {
	dbHost := config.GetConfigValues("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := config.GetConfigValues("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := config.GetConfigValues("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := config.GetConfigValues("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := config.GetConfigValues("DB_NAME")
	if dbName == "" {
		dbName = "transfers"
	}

	dbSSLMode := config.GetConfigValues("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	var gormLogLevel gormlogger.Interface
	if config.GetConfigValues("GORM_LOG_LEVEL") == "debug" {
		gormLogLevel = gormlogger.Default.LogMode(gormlogger.Info)
	} else {
		gormLogLevel = gormlogger.Default.LogMode(gormlogger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogLevel,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		dbLogger.Error(customLogger.LogOptions{
			MethodName: "initializeDB",
			Details:    fmt.Errorf("failed to connect to database: %w", err),
		})
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		dbLogger.Error(customLogger.LogOptions{
			MethodName: "initializeDB",
			Details:    fmt.Errorf("failed to get underlying sql.DB: %w", err),
		})
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	if err := sqlDB.Ping(); err != nil {
		dbLogger.Error(customLogger.LogOptions{
			MethodName: "initializeDB",
			Details:    fmt.Errorf("failed to ping database: %w", err),
		})
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	dbLogger.Debug("Successfully connected to PostgreSQL database using GORM", customLogger.LogOptions{
		MethodName: "initializeDB",
		Details:    "Database connection established",
	})

	return nil
}

func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := DB.AutoMigrate(&models.Account{}, &models.Transaction{}); err != nil {
		dbLogger.Error(customLogger.LogOptions{
			MethodName: "AutoMigrate",
			Details:    fmt.Errorf("failed to auto migrate: %w", err),
		})
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	dbLogger.Debug("Database schema auto-migrated successfully", customLogger.LogOptions{
		MethodName: "AutoMigrate",
		Details:    "All tables created/updated",
	})

	return nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("database not initialized. Call InitDB() first")
	}
	return DB
}

func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}

		if err := sqlDB.Close(); err != nil {
			dbLogger.Error(customLogger.LogOptions{
				MethodName: "CloseDB",
				Details:    fmt.Errorf("failed to close database: %w", err),
			})
			return err
		}

		dbLogger.Debug("Database connection closed", customLogger.LogOptions{
			MethodName: "CloseDB",
			Details:    "Database connection pool closed successfully",
		})
	}
	return nil
}

func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Ping()
}
