package initialize

import (
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/database"
	"go_microservice_backend_api/internal/service"
	"go_microservice_backend_api/internal/service/implement"
)

func InitServiceInterface() {
	queries := database.New(global.Mdb)
	service.InitUserLogin(implement.NewUserLoginImpl(queries))
	service.InitSendMailService(implement.NewSendMailService())
}
