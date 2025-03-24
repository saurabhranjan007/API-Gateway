package unit

import (
	"testing"
	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/validation"

	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger()
	defer logger.SyncLogger()
}

func TestValidateRole(t *testing.T) {

	logger.LogInfo("TestValidateRole", "Test", "Validating role: admin", "")
	err := validation.ValidateRole("admin")
	assert.Nil(t, err)

	logger.LogInfo("TestValidateRole", "Test", "Validating role: invalid_role", "")

	err = validation.ValidateRole("invalid_role")
	logger.LogError("TestValidateRole", "Test", "Invalid role validation", err)

	assert.NotNil(t, err)
	assert.Equal(t, "ROLE invalid_role is not allowed", err.Error())
}

func TestValidateUsername(t *testing.T) {

	logger.LogInfo("TestValidateUsername", "Test", "Validating username: user1", "")
	err := validation.ValidateUsername("user1")
	assert.Nil(t, err)

	logger.LogInfo("TestValidateUsername", "Test", "Validating username: us", "")

	err = validation.ValidateUsername("us")
	logger.LogError("TestValidateUsername", "Test", "Invalid username validation", err)

	assert.NotNil(t, err)
	assert.Equal(t, "username must be at least 4 characters long and contain only letters and numbers", err.Error())
}

func TestValidatePassword(t *testing.T) {

	logger.LogInfo("TestValidatePassword", "Test", "Validating password: Password@123", "")
	err := validation.ValidatePassword("Password@123")
	assert.Nil(t, err)

	logger.LogInfo("TestValidatePassword", "Test", "Validating password: short", "")

	err = validation.ValidatePassword("short")
	logger.LogError("TestValidatePassword", "Test", "Invalid password validation", err)

	assert.NotNil(t, err)
	assert.Equal(t, "password must be at least 8 characters long", err.Error())
}
