package database

import (
	"time"
)

type Auth struct {
	ID           string    `bson:"_id,omitempty"`
	UserID       string    `bson:"user_id"`
	Token        string    `bson:"token"`
	RefreshToken string    `bson:"refresh_token"`
	Created      time.Time `bson:"created"`
}
