package entity

import (
	"crypto/rand"
	"math/big"
	"zeneye-gateway/pkg/logger"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateUserUUID() (string, error) {
	logger.LogInfo("Entity", "generateUserUUID", "Generating user UUID", "")

	const length = 8
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			logger.LogError("Entity", "generateUserUUID", "", err)
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	userUUID := string(result)
	logger.LogInfo("Entity", "generateUserUUID", "User UUID generated", userUUID)
	return userUUID, nil
}
