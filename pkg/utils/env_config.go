package utils

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL               string
	DatabaseName              string
	JWTSecret                 string
	JWTExpiration             int
	RefreshTokenExpiration    int
	RateLimit                 int
	GinMode                   string
	AgentServiceURL           string
	ComplianceServiceURL      string
	ConfigurationServiceURL   string
	NotificationServiceURL    string
	BotDetectionServiceURL    string
	WAFServiceURL             string
	BreachDetectionServiceURL string
	AdminPanelServiceURL      string
	AdminManagementServiceURL string
}

var config *Config
var once sync.Once

func LoadConfig() *Config {
	once.Do(func() {
		env := os.Getenv("GO_ENV")
		if env == "" {
			env = "development"
		}

		envFile := ".env." + env
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Error loading %s file", envFile)
		}

		// Load configurations
		jwtExpiration, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
		if err != nil {
			log.Fatalf("Invalid JWT_EXPIRATION value: %s", os.Getenv("JWT_EXPIRATION"))
		}

		refreshTokenExpiration, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION"))
		if err != nil {
			log.Fatalf("Invalid REFRESH_TOKEN_EXPIRATION value: %s", os.Getenv("REFRESH_TOKEN_EXPIRATION"))
		}

		rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
		if err != nil {
			log.Fatalf("Invalid RATE_LIMIT value: %s", os.Getenv("RATE_LIMIT"))
		}

		config = &Config{
			DatabaseURL:               os.Getenv("DATABASE_URL"),
			DatabaseName:              os.Getenv("DATABASE_NAME"),
			JWTSecret:                 os.Getenv("JWT_SECRET"),
			JWTExpiration:             jwtExpiration,
			RefreshTokenExpiration:    refreshTokenExpiration,
			RateLimit:                 rateLimit,
			GinMode:                   os.Getenv("GIN_MODE"),
			AgentServiceURL:           os.Getenv("AGENT_SERVICE_URL"),
			ComplianceServiceURL:      os.Getenv("COMPLIANCE_SERVICE_URL"),
			ConfigurationServiceURL:   os.Getenv("CONFIGURATION_SERVICE_URL"),
			NotificationServiceURL:    os.Getenv("NOTIFICATION_SERVICE_URL"),
			BotDetectionServiceURL:    os.Getenv("BOT_DETECTION_SERVICE_URL"),
			WAFServiceURL:             os.Getenv("WAF_SERVICE_URL"),
			BreachDetectionServiceURL: os.Getenv("BREACH_DETECTION_SERVICE_URL"),
			AdminPanelServiceURL:      os.Getenv("ADMIN_PANEL_SERVICE_URL"),
			AdminManagementServiceURL: os.Getenv("ADMIN_MANAGEMENT_SERVICE_URL"),
		}

		requiredEnvVars := []string{
			"DATABASE_URL",
			"DATABASE_NAME",
			"JWT_SECRET",
			"JWT_EXPIRATION",
			"REFRESH_TOKEN_EXPIRATION",
			"RATE_LIMIT",
			"GIN_MODE",
			"AGENT_SERVICE_URL",
			"COMPLIANCE_SERVICE_URL",
			"CONFIGURATION_SERVICE_URL",
			"NOTIFICATION_SERVICE_URL",
			"BOT_DETECTION_SERVICE_URL",
			"WAF_SERVICE_URL",
			"BREACH_DETECTION_SERVICE_URL",
			"ADMIN_PANEL_SERVICE_URL",
			"ADMIN_MANAGEMENT_SERVICE_URL",
		}

		for _, envVar := range requiredEnvVars {
			value := os.Getenv(envVar)
			if value == "" {
				log.Fatalf("Environment variable %s is not set", envVar)
			} else {
				log.Printf("Environment variable %s is set to %s", envVar, value)
			}
		}
	})

	return config
}

func GetConfig() *Config {
	if config == nil {
		log.Fatalf("Configuration not loaded. Please call LoadConfig() first.")
	}
	return config
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
