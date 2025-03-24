package loadbalancer

import (
	"errors"
	"net/http"
	"strings"
	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/utils"
)

// RouteRequest routes the request to the appropriate microservice based on the path.
func RouteRequest(req *http.Request) string {

	var serviceURL string
	path := strings.Trim(req.URL.Path, "/")

	switch {
	case strings.HasPrefix(path, "agent"):
		serviceURL = utils.GetEnv("AGENT_SERVICE_URL")
	case strings.HasPrefix(path, "compliance"):
		serviceURL = utils.GetEnv("COMPLIANCE_SERVICE_URL")
	case strings.HasPrefix(path, "config"):
		serviceURL = utils.GetEnv("CONFIGURATION_SERVICE_URL")
	case strings.HasPrefix(path, "notify"):
		serviceURL = utils.GetEnv("NOTIFICATION_SERVICE_URL")
	case strings.HasPrefix(path, "bot-detection"):
		serviceURL = utils.GetEnv("BOT_DETECTION_SERVICE_URL")
	case strings.HasPrefix(path, "waf"):
		serviceURL = utils.GetEnv("WAF_SERVICE_URL")
	case strings.HasPrefix(path, "breach"):
		serviceURL = utils.GetEnv("BREACH_DETECTION_SERVICE_URL")
	case strings.HasPrefix(path, "admin-management"):
		serviceURL = utils.GetEnv("ADMIN_MANAGEMENT_SERVICE_URL")
	default:
		defaultErr := errors.New("default route: " + req.URL.Path)
		logger.LogError("LoadBalancer", "RouteRequest", "Unable to route request", defaultErr)
		return ""
	}

	logger.LogInfo("LoadBalancer", "RouteRequest", "Routed request", map[string]string{"Path": req.URL.Path, "ServiceURL": serviceURL})
	return serviceURL
}
