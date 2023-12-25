package database

import (
	"time"
)

type User struct {
	ID         string    `bson:"_id"`
	TelegramID int64     `bson:"telegram_id"`
	Username   string    `bson:"username"`
	LastAuth   time.Time `bson:"last_auth"`
	Created    time.Time `bson:"created"`
	Updated    time.Time `bson:"updated"`
}
