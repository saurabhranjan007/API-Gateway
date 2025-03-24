package unit

import (
	"testing"
	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/internal/adapter/service"
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/pkg/logger"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&entity.User{}, &entity.RefreshToken{})
	return db
}

func TestCreateUser(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	user := &entity.User{
		Username: "unit_testuser",
		Password: "unit@password",
		Email:    "unit_test@example.com",
		Role:     "admin",
	}

	logger.LogInfo("TestCreateUser", "Test", "Creating user", user)

	err := userService.CreateUser(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, user.UserUUID) // Ensure UserUUID is not empty
}

func TestGetUser(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	user := &entity.User{
		Username: "unit_testuser",
		Password: "unit@password",
		Email:    "unit_test@example.com",
		Role:     "admin",
	}

	logger.LogInfo("TestGetUser", "Test", "Creating user for retrieval", user)

	userService.CreateUser(user)
	result, err := userService.GetUser(user.ID)

	logger.LogInfo("TestGetUser", "Test", "Retrieving user", result)
	assert.Nil(t, err)
	assert.Equal(t, user.Username, result.Username)
}

func TestEditUser(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	user := &entity.User{
		Username: "unit_testuser",
		Password: "unit@password",
		Email:    "unit_test@example.com",
		Role:     "admin",
	}

	logger.LogInfo("TestEditUser", "Test", "Creating user for editing", user)
	userService.CreateUser(user)
	user.Username = "updateduser"
	user.Email = "newemail@example.com"

	logger.LogInfo("TestEditUser", "Test", "Editing user", user)
	err := userService.EditUser(user)
	assert.Nil(t, err)

	result, _ := userService.GetUser(user.ID)
	logger.LogInfo("TestEditUser", "Test", "Retrieving edited user", result)
	assert.Equal(t, "updateduser", result.Username)
	assert.Equal(t, "newemail@example.com", result.Email)
}

func TestDeleteUser(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	user := &entity.User{
		Username: "unit_testuser",
		Password: "unit@password",
		Email:    "unit_test@example.com",
		Role:     "admin",
	}

	logger.LogInfo("TestDeleteUser", "Test", "Creating user for deletion", user)
	userService.CreateUser(user)
	err := userService.DeleteUser(user.ID)

	logger.LogInfo("TestDeleteUser", "Test", "Deleting user", user.ID)
	assert.Nil(t, err)

	result, err := userService.GetUser(user.ID)
	logger.LogInfo("TestDeleteUser", "Test", "Retrieving deleted user", result)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestListUsers(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)

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

	logger.LogInfo("TestListUsers", "Test", "Creating multiple users for listing", user1, user2)
	userService.CreateUser(user1)
	userService.CreateUser(user2)

	users, err := userService.ListUsers()
	logger.LogInfo("TestListUsers", "Test", "Listing all users", users)
	assert.Nil(t, err)
	assert.Len(t, users, 2)
}

func TestCreateSuperadmin(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)

	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}

	logger.LogInfo("TestCreateSuperadmin", "Test", "Creating superadmin", superadmin)
	err := userService.CreateUser(superadmin)
	assert.Nil(t, err)
	assert.NotEmpty(t, superadmin.UserUUID) // Ensure UserUUID is not empty

	// Attempt to create another superadmin
	anotherSuperadmin := &entity.User{
		Username: "superadmin2",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin2@example.com",
		Role:     "superadmin",
	}

	logger.LogInfo("TestCreateSuperadmin", "Test", "Attempting to create another superadmin", anotherSuperadmin)
	err = userService.CreateUser(anotherSuperadmin)
	assert.NotNil(t, err)
	assert.Equal(t, "superadmin already exists", err.Error())
}

func TestIsSuperadminPresent(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)

	logger.LogInfo("TestIsSuperadminPresent", "Test", "Checking if superadmin is present")
	exists, err := userService.IsSuperadminPresent()
	assert.Nil(t, err)
	assert.False(t, exists)

	superadmin := &entity.User{
		Username: "superadmin",
		Password: "SuperSecurePassword@123",
		Email:    "superadmin@example.com",
		Role:     "superadmin",
	}

	logger.LogInfo("TestIsSuperadminPresent", "Test", "Creating superadmin", superadmin)
	err = userService.CreateUser(superadmin)
	assert.Nil(t, err)

	exists, err = userService.IsSuperadminPresent()
	logger.LogInfo("TestIsSuperadminPresent", "Test", "Checking if superadmin is present after creation", exists)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestGenerateRefreshToken(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)

	user := &entity.User{
		Username: "testuser",
		Password: "password@123",
		Email:    "testuser@example.com",
		Role:     "admin",
	}

	logger.LogInfo("TestGenerateRefreshToken", "Test", "Creating user for refresh token generation", user)
	userService.CreateUser(user)

	refreshToken, err := userService.GenerateRefreshToken(user.ID)
	logger.LogInfo("TestGenerateRefreshToken", "Test", "Generating refresh token", refreshToken)
	assert.Nil(t, err)
	assert.NotEmpty(t, refreshToken)
}

func TestRefreshAccessToken(t *testing.T) {
	logger.InitLogger()
	defer logger.SyncLogger()

	db := setupTestDB()
	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)

	user := &entity.User{
		Username: "testuser",
		Password: "password@123",
		Email:    "testuser@example.com",
		Role:     "admin",
	}

	logger.LogInfo("TestRefreshAccessToken", "Test", "Creating user for access token refresh", user)
	userService.CreateUser(user)

	refreshToken, _ := userService.GenerateRefreshToken(user.ID)
	newToken, err := userService.RefreshAccessToken(refreshToken)
	logger.LogInfo("TestRefreshAccessToken", "Test", "Refreshing access token", newToken)
	assert.Nil(t, err)
	assert.NotEmpty(t, newToken)
}
