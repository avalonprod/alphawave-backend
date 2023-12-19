package http

import (
	"net"
	"net/http"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/pkg/limiter"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func corsMiddleware(c *gin.Context) {
	allowedOrigins := []string{
		"https://alphawave.gasstrem.com",
		"http://localhost:3000",
		"http://localhost:7000",
		"http://localhost:3001",
		"http://127.0.0.1:5500",
		"http://localhost",
		"http://localhost:5173",
		"http://localhost:3001",
		"https://plankton-app-kpofc.ondigitalocean.app",
		"plankton-app-kpofc.ondigitalocean.app",
		"https://alpahwave-client.onrender.com",
		"https://oyster-app-4mavy.ondigitalocean.app",
	}

	requestOrigin := c.Request.Header.Get("Origin")
	for _, origin := range allowedOrigins {
		if origin == requestOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			break
		}
	}

	// c.Header("Access-Control-Allow-Origin", "")
	c.Header("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT, PATCH")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func Limit(rps int, burst int, ttl time.Duration) gin.HandlerFunc {
	l := limiter.NewRateLimiter(rps, burst, ttl)

	go l.CleanupVisitors()

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)

		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		if !l.GetVisitor(ip).Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)

			return
		}
		c.Next()
	}
}
