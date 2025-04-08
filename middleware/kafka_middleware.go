package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"

	"intent/config"
)

type CustomResponseRecorder struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func KafkaMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		var jsonData interface{}
		method := c.Request().Method

		var bodyBytes []byte
		var err error

		// Chỉ đọc body nếu là POST, PUT, PATCH
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			bodyBytes, err = io.ReadAll(c.Request().Body)
			if err != nil {
				log.Println("Failed to read request body:", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body")
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
				log.Println("Invalid JSON request:", err)
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
					"error_code": 400,
					"message":    "Invalid JSON format",
					"data":       err.Error(),
				})
			}
		}

		// Ghi response
		rec := &CustomResponseRecorder{
			ResponseWriter: c.Response().Writer,
			Body:           new(bytes.Buffer),
		}
		c.Response().Writer = rec

		err = next(c)

		responseData := rec.Body.String()
		cfg := config.LoadConfig()
		kafkaBroker, topic := config.GetKafkaConfig(cfg)
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{kafkaBroker},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
		defer writer.Close()

		message := map[string]interface{}{
			"request_id": requestID,
			"request":    jsonData,
			"response":   responseData,
			"path":       c.Request().URL.Path,
			"method":     method,
		}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process request")
		}

		if err := writer.WriteMessages(context.Background(), kafka.Message{Value: msg}); err != nil {
			log.Println("Failed to send message:", err)
		} else {
			log.Printf("Sent message with request_id: %s", requestID)
		}

		return err
	}
}
