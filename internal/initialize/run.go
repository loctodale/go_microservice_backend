package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
)

func Run() {
	LoadConfig()
	fmt.Println("Load config mysql", global.Config.Mysql.UserTable.Username)
	InitLogger()
	global.Logger.Info("Config log ok", zap.String("ok", "success"))
	InitReids()
}
