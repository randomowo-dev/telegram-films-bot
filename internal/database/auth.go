package database

import (
	"context"
	"time"

	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	"github.com/randomowo-dev/telegram-films-bot/internal/transport/nosql"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type AuthDB struct {
	db *nosql.DB
}

func (db *AuthDB) Add(ctx context.Context, auth *dbModels.Auth) error {
	auth.Created = time.Now().UTC()
	_, err := db.db.Collection("auth").InsertOne(ctx, auth)
	return err
}

func (db *AuthDB) Exists(ctx context.Context, token string) (bool, error) {
	filter := bson.M{"$or": []bson.M{{"token": token}, {"refresh_token": token}}}
	count, err := db.db.Collection("auth").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (db *AuthDB) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	_, err := db.db.Collection("auth").DeleteMany(ctx, filter)
	return err
}

func NewAuthDB(db *nosql.DB) *AuthDB {
	return &AuthDB{db: db}
}
