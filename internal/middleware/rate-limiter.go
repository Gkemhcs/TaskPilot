package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimiterMiddleware returns a Gin middleware for per-IP rate limiting
func RateLimiterMiddleware(redisClient *redis.Client, logger *logrus.Logger) (gin.HandlerFunc, error) {
	// 10 requests per minute

	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  3,
	}

	// Create Limiter store backed by Redis
	store, err := redisstore.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix:   "limiter",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limiter store: %w", err)
	}

	limiterInstance := limiter.New(store, rate)

	// Middleware logic
	return func(c *gin.Context) {
		// Extract client IP address
		logger.Info("request called")
		clientIP := getClientIP(c.Request)

		ctx := context.Background()
		key := fmt.Sprintf("rate-limit-%s", clientIP)

		// Get rate limit context
		limitCtx, err := limiterInstance.Get(ctx, key)
		if err != nil {
			logger.Info("failed to get rate limit context", "error", err)
			utils.Error(c, http.StatusInternalServerError, "sorry we are experiencing issues, please try again later")

			return
		}
		logger.Info("Rate limit context retrieved", "client_ip", clientIP, "limit", limitCtx.Limit, "remaining", limitCtx.Remaining, "reset", limitCtx.Reset)

		// Add headers for client observability
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limitCtx.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limitCtx.Remaining))
		loc, _ := time.LoadLocation("Asia/Kolkata") // IST location
		resetTime := time.Unix(limitCtx.Reset, 0).In(loc).Format("Mon, 02 Jan 2006 15:04:05 MST")
		c.Header("X-RateLimit-Reset", resetTime)

		

		if limitCtx.Reached {
			logger.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"limit":     limitCtx.Limit,
				"remaining": limitCtx.Remaining,
				"reset":     limitCtx.Reset,
			}).Warn("Rate limit exceeded for client IP")
			utils.Error(c, http.StatusTooManyRequests, "rate limit exceeded. Try again later.")
			return
		}

		c.Next()
	}, nil
}

// Extracts client IP even if behind reverse proxy
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return strings.TrimSpace(ip)
}
