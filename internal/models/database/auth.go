package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Token        string             `bson:"token"`
	RefreshToken string             `bson:"refresh_token"`
	Created      time.Time          `bson:"created"`
}
