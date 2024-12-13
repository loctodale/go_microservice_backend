package private_config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_user/local"
	"go_microservice_backend_api/internal/utils/setPool"
)

type Config struct {
}

func NewUserConfig() *Config {
	return &Config{}
}

func (c *Config) InitUserMysql() {
	m := global.Config.Mysql
	dsnUser := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		m.UserTable.Username, m.UserTable.Password, m.UserTable.Host, m.UserTable.Port, m.UserTable.DbName)
	userDb, err := sql.Open("mysql", dsnUser)
	if err != nil {
		global.Logger.Fatal("Init mysql failed", zap.Error(err))
		panic(err)
	}
	global.Logger.Info("Init mysql success", zap.String("dsn", dsnUser))
	setPool.SetPool(userDb, m.UserTable.MaxIdleConns, m.UserTable.MaxOpenConns, m.UserTable.ConnMaxLifeTime)
	local.UserDb = userDb
}

func (c *Config) InitKafkaProducer() {
	k := &kafka.Writer{
		Addr:     kafka.TCP("kafka0:9092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
	local.UserProducer = k
}
