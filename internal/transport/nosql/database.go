package nosql

import (
	"github.com/randomowo-dev/telegram-films-bot/internal/config"
	"github.com/randomowo-dev/telegram-films-bot/pkg/transport/nosql"
)

type DB struct {
	*nosql.MongoDB
}

func NewDB() *DB {
	return &DB{
		MongoDB: nosql.NewMongoDB(
			nosql.Config{
				Name:    config.DbName,
				Url:     config.DbUrl,
				AppName: config.AppName,
			},
		),
	}
}
