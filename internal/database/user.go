package database

import (
	"context"
	"errors"
	"time"

	dbModels "github.com/randomowo-dev/telegram-films-bot/internal/models/database"
	"github.com/randomowo-dev/telegram-films-bot/internal/transport/nosql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

type UserDB struct {
	db *nosql.DB
}

func (db *UserDB) GetAll(ctx context.Context, offset int64, limit int64) ([]dbModels.User, bool, error) {
	cursor, err := db.db.Collection("user").Find(
		ctx,
		bson.D{},
		options.Find().SetSkip(offset).SetLimit(limit),
	)
	if err != nil {
		return nil, false, err
	}

	var users []dbModels.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, false, err
	}

	count, err := db.db.Collection("user").CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, false, err
	}

	count -= (offset + 1) * limit
	if count < 0 {
		count = 0
	}

	return users, count > 0, nil
}

func (db *UserDB) UpdateByTelegramID(ctx context.Context, telegramID int64, update *dbModels.User) error {
	user, err := db.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	update.Updated = time.Now().UTC()
	_, err = db.db.Collection("user").UpdateByID(ctx, user.ID, bson.D{{"$set", update}})
	if err != nil {
		return err
	}

	update.ID = user.ID

	return nil
}

func (db *UserDB) GetByTelegramID(ctx context.Context, telegramID int64) (*dbModels.User, error) {
	user := new(dbModels.User)
	err := db.db.Collection("user").FindOne(ctx, bson.D{{"telegram_id", telegramID}}).Decode(user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *UserDB) Add(ctx context.Context, user *dbModels.User) error {
	user.Created = time.Now().UTC()
	user.Updated = user.Created
	user.Role = dbModels.UserRole

	res, err := db.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID, _ = res.InsertedID.(primitive.ObjectID)

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
