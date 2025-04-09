package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomResponseRecorder struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func KafkaMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		method := c.Request().Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				log.Println("Failed to read request body:", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body")
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		rec := &CustomResponseRecorder{
			ResponseWriter: c.Response().Writer,
			Body:           new(bytes.Buffer),
		}
		c.Response().Writer = rec
		err := next(c)
		return err
	}
}
