package validation

import (
	"errors"
	"regexp"
	"zeneye-gateway/pkg/logger"
)

var (
	usernameRegex    = regexp.MustCompile(`^[a-zA-Z0-9]{4,}$`)
	specialCharRegex = regexp.MustCompile(`[!@#\$%\^&\*(),.?":{}|<>]`)
	numberRegex      = regexp.MustCompile(`[0-9]`)
	lowerCaseRegex   = regexp.MustCompile(`[a-z]`)
	upperCaseRegex   = regexp.MustCompile(`[A-Z]`)
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func ValidateUsername(username string) error {
	logger.LogInfo("Validation", "ValidateUsername", "Validating username", username)

	if !usernameRegex.MatchString(username) {
		err := errors.New("username must be at least 4 characters long and contain only letters and numbers")
		logger.LogError("Validation", "ValidateUsername", username, err)
		return err
	}

	logger.LogInfo("Validation", "ValidateUsername", "Username is valid", username)
	return nil
}

func ValidatePassword(password string) error {
	logger.LogInfo("Validation", "ValidatePassword", "Validating password", "")

	if len(password) < 8 {
		err := errors.New("password must be at least 8 characters long")
		logger.LogError("Validation", "ValidatePassword", "len(password)", err)
		return err
	}

	if !specialCharRegex.MatchString(password) {
		err := errors.New("password must contain at least one special character")
		logger.LogError("Validation", "ValidatePassword", "!specialCharRegex", err)
		return err
	}

	if !numberRegex.MatchString(password) {
		err := errors.New("password must contain at least one number")
		logger.LogError("Validation", "ValidatePassword", "!numberRegex", err)
		return err
	}

	if !lowerCaseRegex.MatchString(password) {
		err := errors.New("password must contain at least one lowercase letter")
		logger.LogError("Validation", "ValidatePassword", "!lowerCaseRegex", err)
		return err
	}

	if !upperCaseRegex.MatchString(password) {
		err := errors.New("password must contain at least one uppercase letter")
		logger.LogError("Validation", "ValidatePassword", "!upperCaseRegex", err)
		return err
	}

	logger.LogInfo("Validation", "ValidatePassword", "Password is valid", "")
	return nil
}

func ValidateRole(role string) error {
	logger.LogInfo("Validation", "ValidateRole", "Validating role", role)

	validRoles := []string{"superadmin", "admin", "department_admin", "auditor"}
	for _, r := range validRoles {
		if role == r {
			logger.LogInfo("Validation", "ValidateRole", "Role is valid", role)
			return nil
		}
	}

	err := errors.New("ROLE " + role + " is not allowed")
	logger.LogError("Validation", "ValidateRole", role, err)
	return err
}

func ValidateEmail(email string) error {
	logger.LogInfo("Validation", "ValidateEmail", "Validating email", email)

	if !emailRegex.MatchString(email) {
		err := errors.New("invalid email format")
		logger.LogError("Validation", "ValidateEmail", email, err)
		return err
	}

	logger.LogInfo("Validation", "ValidateEmail", "Email is valid", email)
	return nil
}
