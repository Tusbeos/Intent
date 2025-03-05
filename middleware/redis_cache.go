package middleware

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

// Middleware Redis Cache
func RedisCache(redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Thêm Header() để implement http.ResponseWriter
func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Thêm WriteHeader() để implement http.ResponseWriter
func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
