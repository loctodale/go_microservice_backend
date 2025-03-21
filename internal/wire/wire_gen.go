// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"go_microservice_backend_api/internal/controller"
	"go_microservice_backend_api/internal/repo"
	"go_microservice_backend_api/internal/service"
)

// Injectors from user.wire.go:

func InitUserRouterHandle() (*controller.UserController, error) {
	iUserRepository := repo.NewUserRepository()
	iAuthRepository := repo.NewAuthRepository()
	iUserService := service.NewUserService(iUserRepository, iAuthRepository)
	userController := controller.NewUserController(iUserService)
	return userController, nil
}
