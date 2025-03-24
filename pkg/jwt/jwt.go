package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"zeneye-gateway/pkg/logger"
	utils "zeneye-gateway/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte

// InitJWTSecret initializes the JWT secret from environment variables.
func InitJWTSecret() {
	jwtSecret = []byte(utils.GetEnv("JWT_SECRET"))
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	UserUUID string `json:"user_uuid"`
	jwt.RegisteredClaims
}

func getExpirationTime() time.Duration {
	expirationTimeStr := utils.GetEnv("JWT_EXPIRATION")
	expirationTime, err := strconv.Atoi(expirationTimeStr)
	if err != nil {
		logger.LogError("JWT", "getExpirationTime", "Invalid JWT_EXPIRATION, using default 2 hours", err)
		return 2 * time.Hour
	}
	return time.Duration(expirationTime) * time.Hour
}

func getRefreshTokenExpirationTime() time.Duration {
	expirationTimeStr := utils.GetEnv("REFRESH_TOKEN_EXPIRATION")
	expirationTime, err := strconv.Atoi(expirationTimeStr)
	if err != nil {
		logger.LogError("JWT", "getRefreshTokenExpirationTime", "Invalid REFRESH_TOKEN_EXPIRATION, using default 30 days", err)
		return 720 * time.Hour
	}
	return time.Duration(expirationTime) * time.Hour
}

func GenerateToken(userID uint, username, role, userUUID string) (string, error) {
	expirationTime := time.Now().Add(getExpirationTime())
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		UserUUID: userUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.LogError("JWT", "GenerateToken", userID, err)
		return "", err
	}
	logger.LogInfo("JWT", "GenerateToken", "Generated JWT token", signedToken)
	return signedToken, nil
}

func GenerateRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	if _, err := rand.Read(refreshToken); err != nil {
		logger.LogError("JWT", "GenerateRefreshToken", "Error generating refresh token", err)
		return "", err
	}
	refreshTokenString := hex.EncodeToString(refreshToken)
	logger.LogInfo("JWT", "GenerateRefreshToken", "Generated refresh token", refreshTokenString)
	return refreshTokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		logger.LogError("JWT", "ValidateToken", "Error validating token", err)
		return nil, err
	}
	if !token.Valid {
		logger.LogError("JWT", "ValidateToken", "Invalid token signature", err)
		return nil, jwt.ErrSignatureInvalid
	}
	logger.LogInfo("JWT", "ValidateToken", "Token validated successfully", tokenString)
	return claims, nil
}
