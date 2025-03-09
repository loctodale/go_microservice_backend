package service_sendmail

import (
	"go_microservice_backend_api/internal/service_sendmail/consumer"
	"go_microservice_backend_api/internal/service_sendmail/private_config"
	"go_microservice_backend_api/internal/service_sendmail/service"
)

func ServiceSendMailMain() {
	service.InitSendMailServiceInterface()
	private_config.NewSendMailConfig().InitKafkaReader()
	go func() {
		service.InitSendMailServiceInterface()
		for {
			consumer.InitSendMailConsumer().NewSendMailConsumer()
		}
	}()
}
