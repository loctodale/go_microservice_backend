package local

import (
	"database/sql"
	"github.com/segmentio/kafka-go"
)

var (
	ShopDb       *sql.DB
	ShopProducer *kafka.Writer
)
