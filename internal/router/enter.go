package router

import (
	"go_microservice_backend_api/internal/router/manager"
	"go_microservice_backend_api/internal/router/user"
)

type RouterGroup struct {
	User    user.UserRouterGroup
	Manager manager.ManagerRouterGroup
}

var RouterGroupApp = new(RouterGroup)
