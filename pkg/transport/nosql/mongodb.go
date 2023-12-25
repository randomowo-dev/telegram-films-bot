package nosql

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBInterface interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Collection(string)
}

type MongoDB struct {
	conn    *mongo.Client
	options *options.ClientOptions
	dbName  string
}

func (db *MongoDB) Connect(ctx context.Context) error {
	if db.conn != nil {
		return db.conn.Ping(ctx, nil)
	}

	var err error
	db.conn, err = mongo.Connect(ctx, db.options)
	if err != nil {
		return err
	}

	return db.conn.Ping(ctx, nil)
}

func (db *MongoDB) Disconnect(ctx context.Context) error {
	if db.conn == nil {
		return nil
	}

	defer func() {
		db.conn = nil
	}()

	return db.conn.Disconnect(ctx)
}

func (db *MongoDB) Collection(collection string) *mongo.Collection {
	if db.conn == nil {
		panic(fmt.Errorf("no connection to db"))
	}
	return db.conn.Database(db.dbName).Collection(collection)
}

type Config struct {
	AppName string
	Url     string
	Name    string
}

func NewMongoDB(config Config) *MongoDB {
	optionsClient := options.Client()
	optionsClient.SetAppName(config.AppName)
	optionsClient.ApplyURI(config.Url)

	return &MongoDB{
		conn:    nil,
		options: optionsClient,
		dbName:  config.Name,
	}
}
