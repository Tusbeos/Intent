package middleware

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// Middleware Redis Cache
func RedisCache(redisClient *redis.Client, cacheableRoutes map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Chỉ cache request GET
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			// Kiểm tra xem route này có trong danh sách cần cache không
			if _, shouldCache := cacheableRoutes[c.Path()]; !shouldCache {
				return next(c)
			}

			ctx := context.Background()
			key := "cache:" + c.Request().RequestURI

			// Kiểm tra response trong Redis
			cachedResponse, err := redisClient.Get(ctx, key).Result()
			if err == nil {
				// Nếu có cache, trả về response từ Redis
				return c.JSONBlob(http.StatusOK, []byte(cachedResponse))
			}

			// Nếu không có cache, tiếp tục request
			resBody := new(bytes.Buffer)
			writer := &responseWriter{
				ResponseWriter: c.Response().Writer,
				body:           resBody,
			}
			c.Response().Writer = writer

			// Gọi API gốc
			if err := next(c); err != nil {
				return err
			}

			// Lưu response vào Redis với TTL 60 giây
			redisClient.Set(ctx, key, resBody.String(), 60*time.Second)

			return nil
		}
	}
}

// Custom Response Writer để lưu response
type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
