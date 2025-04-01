package message_queue

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaWorker struct {
	broker    string
	topic     string
	processor *Processor
}

func NewKafkaWorker(broker, topic string, processor *Processor) *KafkaWorker {
	return &KafkaWorker{
		broker:    broker,
		topic:     topic,
		processor: processor,
	}
}

func (w *KafkaWorker) Start() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{w.broker},
		Topic:   w.topic,
		GroupID: "user-actions",
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("[Kafka] Error reading message:", err)
			continue
		}

		// Gửi message tới Processor để xử lý
		go w.processor.ProcessMessage(msg.Value)
	}
}
