package local

import (
	"database/sql"
	"github.com/segmentio/kafka-go"
)

var (
	ProductDb       *sql.DB
	ProductProducer *kafka.Writer
)
