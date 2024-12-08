package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/consumer"
)

func Run() *gin.Engine {
	LoadConfig()
	fmt.Println("Load config mysql", global.Config.Mysql.Username)
	InitLogger()
	global.Logger.Info("Config log ok", zap.String("ok", "success"))
	InitMysql()
	InitServiceInterface()
	InitReids()
	InitKafka()

	r := InitRouter()
	go func() {
		consumer := consumer.SendMailConsumer{}
		consumer.NewSendMailConsumer()
	}()
	return r
}
