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

			// Kiểm tra số request hiện tại
			countStr, err := redisClient.Get(ctx, key).Result()
			count, _ := strconv.Atoi(countStr)

			if err == nil && count >= limit {
				// Nếu quá giới hạn, trả về lỗi quá nhiều request
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests, please try again later.",
				})
			}

			// Nếu chưa đến giới hạn, tăng số request
			redisClient.Incr(ctx, key)
			redisClient.Expire(ctx, key, duration)

			return next(c)
		}
	}
}
