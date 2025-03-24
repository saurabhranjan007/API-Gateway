package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	internal "zeneye-gateway/internal/adapter/http"
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/pkg/logger"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUserIntegration(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestCreateUserIntegration", "Test", "Starting integration test for creating a user", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	user := map[string]string{
		"username": "testuser",
		"password": "password",
		"email":    "test@example.com",
		"role":     "admin",
	}
	jsonValue, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	logger.LogInfo("TestCreateUserIntegration", "Test", "Response", w)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser entity.User
	db.Where("username = ?", "testuser").First(&createdUser)
	assert.NotEmpty(t, createdUser.UserUUID)
	assert.Equal(t, 8, len(createdUser.UserUUID))
}

func TestEditUserIntegration(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestEditUserIntegration", "Test", "Starting integration test for editing a user", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	user := entity.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Role:     "admin",
	}
	db.Create(&user)

	updatedUser := map[string]string{
		"username": "updateduser",
		"email":    "newemail@example.com",
	}
	jsonValue, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/"+strconv.Itoa(int(user.ID)), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	logger.LogInfo("TestEditUserIntegration", "Test", "Response", w)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedUserFromDB entity.User
	db.First(&updatedUserFromDB, user.ID)
	assert.Equal(t, "updateduser", updatedUserFromDB.Username)
	assert.Equal(t, "newemail@example.com", updatedUserFromDB.Email)
	assert.Equal(t, user.UserUUID, updatedUserFromDB.UserUUID) // Ensure user_uuid is not changed
}

func TestDeleteUserIntegration(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestDeleteUserIntegration", "Test", "Starting integration test for deleting a user", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	user := entity.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Role:     "admin",
	}
	db.Create(&user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(int(user.ID)), nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	logger.LogInfo("TestDeleteUserIntegration", "Test", "Response", w)

	assert.Equal(t, http.StatusOK, w.Code)

	var deletedUser entity.User
	result := db.First(&deletedUser, user.ID)
	assert.NotNil(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestGetUserIntegration(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestGetUserIntegration", "Test", "Starting integration test for getting a user", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	user := entity.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Role:     "admin",
	}
	db.Create(&user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(user.ID)), nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	logger.LogInfo("TestGetUserIntegration", "Test", "Response", w)

	assert.Equal(t, http.StatusOK, w.Code)

	var result entity.User
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.UserUUID, result.UserUUID)
}

func TestListUsersIntegration(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestListUsersIntegration", "Test", "Starting integration test for listing users", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	user1 := &entity.User{
		Username: "testuser1",
		Password: "password1",
		Email:    "test1@example.com",
		Role:     "admin",
	}
	user2 := &entity.User{
		Username: "testuser2",
		Password: "password2",
		Email:    "test2@example.com",
		Role:     "admin",
	}
	db.Create(&user1)
	db.Create(&user2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	logger.LogInfo("TestListUsersIntegration", "Test", "Response", w)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []entity.User
	json.Unmarshal(w.Body.Bytes(), &users)
	assert.Len(t, users, 2)
}

func TestCheckSuperadmin(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestCheckSuperadmin", "Test", "Starting integration test for checking superadmin", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/superadmin/check", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	logger.LogInfo("TestCheckSuperadmin", "Test", "Response", w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"superadmin_exists":true`)
}

func TestCreateSuperadmin(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := SetupTestDB()
	router := internal.SetupRouter(db)

	logger.LogInfo("TestCreateSuperadmin", "Test", "Starting integration test for creating superadmin", "")

	// Create a superadmin user to get a valid token
	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}
	db.Create(&superadmin)
	token := GenerateTestToken(superadmin.ID)

	// Attempt to create another superadmin
	w := httptest.NewRecorder()
	reqBody := `{
        "username": "superadmin2",
        "password": "SuperSecurePassword@123",
        "email": "superadmin2@example.com",
        "role": "superadmin"
    }`
	req, _ := http.NewRequest("POST", "/superadmin/create", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Body = ioutil.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)

	logger.LogInfo("TestCreateSuperadmin", "Test", "Response", w)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Superadmin already exists")
}
