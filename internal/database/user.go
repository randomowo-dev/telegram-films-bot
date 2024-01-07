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

func (db *UserDB) UpdateByTelegramID(ctx context.Context, telegramID int64) (string, error) {
	user, err := db.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}

	now := time.Now().UTC()
	_, err = db.db.Collection("user").UpdateByID(
		ctx, user.ID, bson.D{
			{
				"$set",
				bson.D{
					{"last_auth", now},
					{"updated", now},
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (db *UserDB) GetByTelegramID(ctx context.Context, telegramID int64) (*dbModels.User, error) {
	user := new(dbModels.User)
	err := db.db.Collection("user").FindOne(ctx, bson.D{{"telegram_id", telegramID}}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *UserDB) Add(ctx context.Context, user *dbModels.User) error {
	user.Created = time.Now().UTC()
	user.Updated = user.Created
	res, err := db.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID, _ = res.InsertedID.(string)

	return nil
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
