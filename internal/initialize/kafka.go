package initialize

import (
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/global"
	"log"
)

// Init kafka Producer
var kafkaProducer *kafka.Writer

func InitKafka() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatal(err.Error())
	}
}
