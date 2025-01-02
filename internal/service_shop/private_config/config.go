package private_config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_shop/local"
	"go_microservice_backend_api/internal/utils/setPool"
)

type Config struct {
}

func NewShopConfig() *Config {
	return &Config{}
}

func (c *Config) InitShopMysql() {
	m := global.Config.Mysql
	dsnShop := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		m.ShopTable.Username, m.ShopTable.Password, m.ShopTable.Host, m.ShopTable.Port, m.ShopTable.DbName)
	shopDb, err := sql.Open("mysql", dsnShop)
	if err != nil {
		global.Logger.Fatal("Init mysql failed", zap.Error(err))
		panic(err)
	}
	global.Logger.Info("Init mysql success", zap.String("dsn", dsnShop))
	setPool.SetPool(shopDb, m.ShopTable.MaxIdleConns, m.ShopTable.MaxOpenConns, m.ShopTable.ConnMaxLifeTime)
	local.ShopDb = shopDb
}

func (c *Config) InitShopKafkaProducer() {
	k := &kafka.Writer{
		Addr:     kafka.TCP("kafka0:9092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
	local.ShopProducer = k
}
