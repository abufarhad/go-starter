package middlewares

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/abufarhad/golang-starter-rest-api/internal/config"
	"strings"
	"time"
)

func Attach(r *gin.Engine) {

	// Gin middlewares
	r.Use(gin.LoggerWithFormatter(GinLogFormatter))
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(JWTWithConfig(JWTConfig{
		Skipper: func(c *gin.Context) bool {
			if strings.HasPrefix(c.Request.URL.Path, "/api/v1/docs") {
				return true
			}
			switch c.Request.URL.Path {
			case
				"/api/v1",
				"/api/v1/h34l7h",
				"/api/v1/user/signup",
				"/api/v1/token/refresh",
				"/api/v1/login":
				return true
			default:
				return false
			}
		},
		SigningKey: []byte(config.Jwt().AccessTokenSecret),
		ContextKey: config.Jwt().ContextKey,
	}))
}

// GinLogFormatter is a custom log formatter for Gin
func GinLogFormatter(params gin.LogFormatterParams) string {
	// Your custom log format here
	return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s  %s\n%s",
		params.TimeStamp.Format(time.RFC3339),
		params.StatusCode,
		params.Latency,
		params.ClientIP,
		params.Method,
		params.Path,
		params.ErrorMessage,
	)
}
