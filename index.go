package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexDetail struct {
	V                  int32          `bson:"v"`
	Key                map[string]int `bson:"key"`
	NS                 string         `bson:"ns"`
	Name               string         `bson:"name"`
	Unique             bool           `bson:"unique"`
	ExpireAfterSeconds int32          `bson:"expireAfterSeconds"`
	Sparse             bool           `bson:"sparse"`
	Background         bool           `bson:"background"`
}

func (c *Command) IndexList() (result []IndexDetail, err error) {
	cursor, err := c.collection.Indexes().List(context.Background())
	if err != nil {
		return
	}
	err = cursor.All(context.Background(), &result)

	return
}

// 建立index
type InsertIndexModel struct {
	Keys    interface{}
	Options *options.IndexOptions
}

type IndexOptions = options.IndexOptions

func NewIndexOption() *IndexOptions {
	return &IndexOptions{}
}

func (c *Command) CreateIndex(model InsertIndexModel) (indexName string, err error) {
	opts := options.CreateIndexes()
	m := mongo.IndexModel{
		Keys:    model.Keys,
		Options: model.Options,
	}
	indexName, err = c.collection.Indexes().CreateOne(context.Background(), m, opts)
	return
}

func (c *Command) CreateIndexes(models ...InsertIndexModel) (indexNames []string, err error) {
	if len(models) == 0 {
		return
	}
	opts := options.CreateIndexes()
	indexs := []mongo.IndexModel{}
	for _, index := range models {
		indexs = append(indexs, mongo.IndexModel{
			Keys:    index.Keys,
			Options: index.Options,
		})
	}

	indexNames, err = c.collection.Indexes().CreateMany(context.Background(), indexs, opts)
	return
}

func (c *Command) DropIndex(name string) (err error) {
	_, err = c.collection.Indexes().DropOne(context.Background(), name)
	return
}

func (c *Command) DropAllIndex() (err error) {
	_, err = c.collection.Indexes().DropAll(context.Background())
	return
}
