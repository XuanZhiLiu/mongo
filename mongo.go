package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *MongoDB
	// db         *mongo.Database
	defaultCtx = context.Background()
)

var client *mongo.Client

type CollectionBase struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

type MongoDB struct {
	db      *mongo.Database
	options *options.ClientOptions
	auth    *options.Credential
	logMode logModeValue
	logger  logger
}

type logModeValue int

const (
	defaultLogMode logModeValue = iota
	noLogMode
)

func Init(dbHost, dbUser, dbPassword, dbName string, opts ...MongoOption) (err error) {
	db, err = initMongoDB(dbHost, dbUser, dbPassword, dbName, opts...)
	return
}

func InitMongoDatabase(dbHost, dbUser, dbPassword, dbName string, opts ...MongoOption) (database *MongoDB, err error) {
	if client != nil {
		database = &MongoDB{db: client.Database(dbName)}
		return
	}

	database, err = initMongoDB(dbHost, dbUser, dbPassword, dbName, opts...)
	return
}

func initMongoDB(dbHost, dbUser, dbPassword, dbName string, opts ...MongoOption) (mongoClient *MongoDB, err error) {
	mongoClient = &MongoDB{}
	mongoClient.options = options.Client().ApplyURI("mongodb://" + dbHost).SetMaxConnIdleTime(5 * time.Second).SetMaxPoolSize(200)
	if dbUser != "" {
		mongoClient.auth = &options.Credential{
			Username: dbUser,
			Password: dbPassword,
		}
	}

	for _, opt := range opts {
		err = opt(mongoClient)
		if err != nil {
			return
		}
	}

	if mongoClient.auth != nil {
		mongoClient.options.SetAuth(*mongoClient.auth)
	}

	if client, err = mongo.Connect(context.Background(), mongoClient.options); err != nil {
		println(err.Error())
		return
	}
	mongoClient.db = client.Database(dbName)
	return
}
