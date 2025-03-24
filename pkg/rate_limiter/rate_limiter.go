package rate_limiter

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"zeneye-gateway/pkg/logger"

	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex
var rateLimit int

func init() {
	logger.InitLogger() // ensure log

	// limit, err := strconv.Atoi(utils.GetEnv("RATE_LIMIT"))
	limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		limit = 100 // Default limit
		logger.LogError("RateLimiter", "init", "Invalid RATE_LIMIT environment variable; using default rate limit", err)
	} else {
		logger.LogInfo("RateLimiter", "init", "Rate limit initialized from environment variable", map[string]int{"RateLimit": limit})
	}
	rateLimit = limit
}

func getVisitor(ip string) *rate.Limiter {

	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rateLimit), rateLimit)
		visitors[ip] = limiter
		logger.LogInfo("RateLimiter", "getVisitor", "Created new rate limiter for IP", map[string]string{"IP": ip})
	}

	return limiter
}

func AllowRequest(ip string) bool {

	limiter := getVisitor(ip)
	allowed := limiter.Allow()

	if !allowed {
		notAllowedErr := errors.New("REQUEST NOT ALLOWED DUE TO RATE LIMITER")
		logger.LogError("RateLimiter", "AllowRequest", map[string]string{"IP": ip}, notAllowedErr)
	}

	return allowed
}
