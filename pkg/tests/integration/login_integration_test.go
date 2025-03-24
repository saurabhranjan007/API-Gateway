package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	internal "zeneye-gateway/internal/adapter/http"
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/pkg/logger"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {

	logger.InitLogger() // ensure logger
	defer logger.SyncLogger()

	logger.LogInfo("TestLogin", "TestLogin", "Starting TestLogin", "")

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	// Create a user with hashed password
	password := "password@123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogFatal("TestLogin", "HashPassword", "", err)
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := &entity.User{
		Username: "user1",
		Password: string(hashedPassword),
		Email:    "user1@example.com",
		Role:     "admin",
	}

	if err := db.Create(user).Error; err != nil {
		logger.LogFatal("TestLogin", "CreateUser", user, err)
		t.Fatalf("Failed to create user: %v", err)
	}

	logger.LogInfo("TestLogin", "CreateUser", "User created successfully", user)

	// Test login
	reqBody := `{
        "username": "user1",
        "password": "password@123"
    }`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	logger.LogInfo("TestLogin", "LoginRequest", "Login request sent", reqBody)
	logger.LogInfo("TestLogin", "LoginResponse", "Login response received", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200, got %d", w.Code)
	assert.NotEmpty(t, w.Header().Get("Authorization"), "Authorization header should not be empty")
	assert.NotEmpty(t, w.Header().Get("Refresh-Token"), "Refresh-Token header should not be empty")

	logger.LogInfo("TestLogin", "TestLogin", "TestLogin completed successfully")
}
