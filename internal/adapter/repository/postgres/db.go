package postgres

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"zeneye-gateway/internal/adapter/service"
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB initializes the database connection
func NewDB() (*gorm.DB, error) {
	dsn := utils.GetEnv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.LogError("NewDB", "Failed to open database", "", err)
		return nil, err
	}
	return db, nil
}

// getDSNWithoutDB returns the DSN without the database name
func getDSNWithoutDB(dsn string) (string, string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		logger.LogError("getDSNWithoutDB", "Failed to parse DSN", dsn, err)
		return "", "", err
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	u.Path = "/postgres"
	return u.String(), dbName, nil
}

// checkAndCreateDatabase checks if the database exists and creates it if it doesn't
func checkAndCreateDatabase() error {
	dsn := utils.GetEnv("DATABASE_URL")
	dsnWithoutDB, dbName, err := getDSNWithoutDB(dsn)
	if err != nil {
		logger.LogError("checkAndCreateDatabase", "Failed to parse DSN", dsn, err)
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Connect to the PostgreSQL server
	db, err := sql.Open("postgres", dsnWithoutDB)
	if err != nil {
		logger.LogError("checkAndCreateDatabase", "Failed to connect to PostgreSQL server", dsnWithoutDB, err)
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}
	defer db.Close()

	// Check if the database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", dbName)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		logger.LogError("checkAndCreateDatabase", "Failed to check if database exists", query, err)
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create the database if it doesn't exist
	if !exists {
		query = fmt.Sprintf("CREATE DATABASE %s", dbName)
		_, err = db.Exec(query)
		if err != nil {
			logger.LogError("checkAndCreateDatabase", "Failed to create database", query, err)
			return fmt.Errorf("failed to create database: %w", err)
		}
		logger.LogInfo("checkAndCreateDatabase", "Database created successfully", dbName)
	}

	return nil
}

// migrateDB runs the migrations using golang-migrate
func migrateDB() {
	logger.LogInfo("repository/db", "migrateDB", "Running Migration DB", "")

	dsn := utils.GetEnv("DATABASE_URL")
	migrationsDir, _ := filepath.Abs("db/migrations")
	m, err := migrate.New(
		"file://"+migrationsDir,
		dsn,
	)
	if err != nil {
		logger.LogFatal("migrateDB", "Could not start migration", "", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.LogFatal("migrateDB", "Could not run migration", "", err)
	}
	logger.LogInfo("repository/db", "migrateDB", "Migrations ran successfully", "")
}

// InitDB initializes the database and runs migrations
func InitDB() *gorm.DB {
	if err := checkAndCreateDatabase(); err != nil {
		logger.LogFatal("InitDB", "Failed to check and create database", "", err)
	}

	migrateDB()
	db, err := NewDB()
	if err != nil {
		logger.LogFatal("InitDB", "Failed to initialize database", "", err)
	}

	// Auto migrate the schemas
	db.AutoMigrate(&entity.User{}, &entity.Session{}, &entity.RefreshToken{})

	// Check if superadmin exists, and log the result
	repo := NewUserRepository(db)
	userService := service.NewUserService(repo)
	superadminExists, err := userService.IsSuperadminPresent()
	if err != nil {
		logger.LogFatal("InitDB", "Error checking for superadmin", "", err)
	}

	if superadminExists {
		logger.LogInfo("InitDB", "Superadmin exists in the database", "", superadminExists)
	} else {
		logger.LogInfo("InitDB", "Superadmin does not exist in the database", "", superadminExists)
	}

	return db
}
