package initialize

import (
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/global"
	"log"
	"time"
)

// Init kafka Producer
var kafkaProducer *kafka.Writer

func InitKafka() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("kafka0:9092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
	//kafka.NewR{Brokers: []string{"localhost:29092"}, GroupID: "group-verify-otp", Topic: "otp-auth-topic"},
	global.KafkaConsumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"kafka0:9092"},
		GroupID:        "group-verify-otp",
		Topic:          "otp-auth-topic",
		CommitInterval: time.Second,
	})
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatal(err.Error())
	}

	if err := global.KafkaConsumer.Close(); err != nil {
		log.Fatal(err.Error())
	}
}
