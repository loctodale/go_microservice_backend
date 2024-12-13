package local

import (
	"database/sql"
	"github.com/segmentio/kafka-go"
)

var (
	UserDb       *sql.DB
	UserProducer *kafka.Writer
)
