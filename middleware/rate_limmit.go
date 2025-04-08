package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware giới hạn số request theo IP
func RateLimitMiddleware(redisClient *redis.Client, limit int, duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()
			ip := c.RealIP()
			key := "rate_limit:" + ip
			countStr, err := redisClient.Get(ctx, key).Result()
			count, _ := strconv.Atoi(countStr)

			if err == nil && count >= limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests, please try again later.",
				})
			}
			redisClient.Incr(ctx, key)
			redisClient.Expire(ctx, key, duration)

			return next(c)
		}
	}
}
