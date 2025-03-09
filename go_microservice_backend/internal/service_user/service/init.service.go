package service

import (
	"go_microservice_backend_api/internal/service_user/database"
	"go_microservice_backend_api/internal/service_user/local"
	implement2 "go_microservice_backend_api/internal/service_user/service/implement"
)

func InitUserServiceInterface() {
	userQueries := database.New(local.UserDb)
	InitUserLogin(implement2.NewUserLoginImpl(userQueries))
}
