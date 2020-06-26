package mongo

import (
	"context"

	"github.com/XuanZhiLiu/mongo/bsontool"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type CollectionDetail struct {
	Name    string           `bson:"name"`
	Type    string           `bson:"type"`
	Options CollectionOption `bson:"options"`
	Info    CollectionInfo   `bson:"info"`
	Index   CollectionIndex  `bson:"idIndex"`
}

type CollectionInfo struct {
	ReadOnly bool             `bson:"readOnly"`
	UUID     primitive.Binary `bson:"uuid"`
}

type CollectionOption struct {
	Capped bool  `bson:"capped"`
	Size   int32 `bson:"size"`
}

type CollectionIndex struct {
	V    int32     `bson:"v"`
	Key  bsonx.Doc `bson:"key"`
	Name string    `bson:"name"`
	NS   string    `bson:"ns"`
}

func ListCollections(filter interface{}) (result []*CollectionDetail, err error) {
	return db.ListCollections(filter)
}

func (d *MongoDB) ListCollections(filter interface{}) (result []*CollectionDetail, err error) {
	if filter == nil {
		filter = bsontool.NewBsonD().NotEqual("name", true, "system.profile")
	}
	var cur *mongo.Cursor
	cur, err = d.db.ListCollections(context.Background(), filter)
	if err != nil {
		return
	}
	err = cur.All(context.Background(), &result)
	return
}
