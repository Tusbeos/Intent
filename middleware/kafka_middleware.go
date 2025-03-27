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

// KafkaMiddleware gửi request vào Kafka với request_id
func KafkaMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Tạo request_id (UUID)
		requestID := uuid.New().String()
		c.Set("request_id", requestID) // Lưu vào context Echo

		// Đọc body request và khôi phục lại
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Println("[KafkaMiddleware] Failed to read request body:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read request body")
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Kiểm tra JSON có phải object (`{}`) hoặc array (`[]`) không
		var jsonData interface{}
		if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
			log.Println("[KafkaMiddleware] Invalid JSON request:", err)
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
				"error_code": 400,
				"message":    "Invalid JSON format",
				"data":       err.Error(),
			})
		}

		// Ghi nhận response
		rec := &CustomResponseRecorder{
			ResponseWriter: c.Response().Writer,
			Body:           new(bytes.Buffer),
		}
		c.Response().Writer = rec

		// Gọi API
		err = next(c)

		// Đọc response body
		responseData := rec.Body.String()

		// Gửi vào Kafka
		cfg := config.LoadConfig()
		kafkaBroker, topic := config.GetKafkaConfig(cfg)
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{kafkaBroker},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
		defer writer.Close()

		// Gói dữ liệu gửi đi
		message := map[string]interface{}{
			"request_id": requestID,
			"request":    jsonData, // Lưu nguyên dạng JSON
			"response":   responseData,
			"path":       c.Request().URL.Path,
			"method":     c.Request().Method,
		}
		msg, err := json.Marshal(message)
		if err != nil {
			log.Println("[KafkaMiddleware] Failed to marshal message:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process request")
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Value: msg,
		})
		if err != nil {
			log.Println("[Kafka] Failed to send message:", err)
		} else {
			log.Printf("[Kafka] Sent message with request_id: %s", requestID)
		}

		return err
	}
}
