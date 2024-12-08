//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"go_microservice_backend_api/internal/controller"
	"go_microservice_backend_api/internal/repo"
	"go_microservice_backend_api/internal/service"
)

func InitUserRouterHandle() (*controller.UserController, error) {
	wire.Build(
		repo.NewUserRepository,
		repo.NewAuthRepository,
		service.NewUserService,
		controller.NewUserController,
	)
	return new(controller.UserController), nil
}
