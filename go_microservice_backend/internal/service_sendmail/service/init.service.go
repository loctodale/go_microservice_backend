package service

import "go_microservice_backend_api/internal/service_sendmail/service/implement"

func InitSendMailServiceInterface() {
	InitSendMailService(implement.NewSendMailService())
}
