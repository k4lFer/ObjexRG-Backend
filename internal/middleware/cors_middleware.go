package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var defaultAllowedOrigins = []string{
	"http://localhost:5173",
	"http://127.0.0.1:5173",
	"http://localhost:4173",
	"http://127.0.0.1:4173",
	"http://localhost:3000",
	"http://127.0.0.1:5500",
}

func allowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return defaultAllowedOrigins
	}

	origins := make([]string, 0)
	for _, origin := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	if len(origins) == 0 {
		return defaultAllowedOrigins
	}

	return origins
}

func isAllowedOrigin(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, candidate := range allowed {
		if candidate == "*" || strings.EqualFold(candidate, origin) {
			return true
		}
	}
	return false
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := allowedOrigins()

		if isAllowedOrigin(origin, allowed) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
