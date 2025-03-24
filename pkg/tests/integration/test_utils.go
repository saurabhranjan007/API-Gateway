package integration

import (
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/pkg/jwt"
	"zeneye-gateway/pkg/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	logger.LogInfo("SetupTestDB", "SetupTestDB", "Setting up the test database", "")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		logger.LogFatal("SetupTestDB", "OpenDatabase", "", err)
		panic("failed to connect database")
	}

	logger.LogInfo("SetupTestDB", "OpenDatabase", "Database connection established", "")

	err = db.AutoMigrate(&entity.User{}, &entity.RefreshToken{})
	if err != nil {
		logger.LogFatal("SetupTestDB", "AutoMigrate", "", err)
		panic("failed to migrate database schema")
	}

	logger.LogInfo("SetupTestDB", "AutoMigrate", "Database schema migrated successfully", "")

	return db
}

func GenerateTestToken(userID uint) string {
	logger.LogInfo("GenerateTestToken", "GenerateTestToken", "Generating test token", userID)

	token, err := jwt.GenerateToken(userID)
	if err != nil {
		logger.LogFatal("GenerateTestToken", "GenerateToken", userID, err)
		panic("failed to generate test token")
	}

	logger.LogInfo("GenerateTestToken", "GenerateToken", "Test token generated successfully", token)

	return "Bearer " + token
}
