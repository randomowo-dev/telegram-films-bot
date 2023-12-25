package database

import (
	"context"
	"time"

	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	"github.com/randomowo-dev/telegram-films-bot/internal/transport/nosql"
	"go.mongodb.org/mongo-driver/bson"
)

type UserDB struct {
	db *nosql.DB
}

func (db *UserDB) UpdateByTelegramID(ctx context.Context, user *dbModels.User) error {
	user.Updated = time.Now()
	res, err := db.db.Collection("user").UpdateOne(
		ctx,
		bson.D{{"telegram_id", user.TelegramID}},
		bson.D{{"$set", user}},
	)
	if err != nil {
		return err
	}

	if res.ModifiedCount > 0 {
		return nil
	}

	return db.Add(ctx, user)
}

func (db *UserDB) FindByTelegramID(ctx context.Context, telegramID int64) (*dbModels.User, error) {
	user := new(dbModels.User)
	err := db.db.Collection("user").FindOne(ctx, bson.D{{"telegram_id", telegramID}}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *UserDB) Add(ctx context.Context, user *dbModels.User) error {
	user.Created = time.Now()
	user.Updated = time.Now()
	_, err := db.db.Collection("user").InsertOne(ctx, user)
	return err
}

func (db *UserDB) DeleteByID(ctx context.Context, id string) error {
	_, err := db.db.Collection("user").DeleteOne(ctx, bson.D{{"_id", id}})
	return err
}

func (db *UserDB) DeleteByTelegramID(ctx context.Context, telegramID int64) error {
	_, err := db.db.Collection("user").DeleteOne(ctx, bson.D{{"telegram_id", telegramID}})
	return err
}

func NewUserDB(db *nosql.DB) *UserDB {
	return &UserDB{db: db}
}
