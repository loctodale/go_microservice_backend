package private_config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_product/local"
)

type Config struct{}

func NewProductConfig() *Config {
	return &Config{}
}

func (c *Config) InitProductSql() {
	m := global.Config.Mysql.ProductTable
	dsnProduct := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		m.Username, m.Password, m.Host, m.Port, m.DbName)
	shopDb, err := sql.Open("mysql", dsnProduct)
	if err != nil {
		global.Logger.Error("Init Mysql Shop Service Failed", zap.Error(err))
		panic(err)
	}
	local.ProductDb = shopDb
}

func (c *Config) InitKafkaProducer() {
	k := &kafka.Writer{
		Addr:     kafka.TCP("kafka0:9092"),
		Topic:    "product-topic",
		Balancer: &kafka.LeastBytes{},
	}
	local.ProductProducer = k
}
