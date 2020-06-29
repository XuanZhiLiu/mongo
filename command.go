package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/georgeliu825/mongo/bsontool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Command struct {
	Context     context.Context
	isDebug     bool
	collection  *mongo.Collection
	sort        bsonx.Doc
	skip        *int64
	limit       *int64
	isPaging    bool
	projection  bson.D
	condition   bsontool.D
	pipeline    Pipeline
	updateModel UpdateModel
}

type UpdateModel struct {
	updateOperator bson.M
}

type CmdResult struct {
	ID    primitive.ObjectID
	IDs   []primitive.ObjectID
	Count int64
}

type UpdateParam struct {
	FieldName string
	Document  interface{}
}

func NewCommand(ctx context.Context, collection string) *Command {
	return db.NewCommand(ctx, collection)
}

func (d *MongoDB) NewCommand(ctx context.Context, collection string) *Command {
	c := &Command{}
	c.Context = ctx
	c.collection = d.db.Collection(collection)
	c.condition = bsontool.NewBsonD()
	c.pipeline = NewPipeline()
	c.setPaging(ctx)
	return c
}

func (c *Command) Debug() *Command {
	c.isDebug = true
	return c
}

func (c *Command) DropCollection() (err error) {
	err = c.collection.Drop(context.Background())
	return
}

func (c *Command) Insert(document interface{}) (result CmdResult, err error) {
	return c.insert(document)
}

func (c *Command) InsertMany(documents ...interface{}) (result CmdResult, err error) {
	return c.insertMany(documents)
}

func (c *Command) insert(documents ...interface{}) (result CmdResult, err error) {
	var document interface{}
	if len(documents) > 0 {
		document = documents[0]
	}
	var res *mongo.InsertOneResult
	res, err = c.collection.InsertOne(context.Background(), document)
	if err != nil || res == nil {
		return
	}
	result = CmdResult{ID: res.InsertedID.(primitive.ObjectID)}
	return
}

func (c *Command) insertMany(documents []interface{}) (result CmdResult, err error) {
	var res *mongo.InsertManyResult
	res, err = c.collection.InsertMany(context.Background(), documents)
	if err != nil || res == nil {
		return
	}
	for _, id := range res.InsertedIDs {
		result.IDs = append(result.IDs, id.(primitive.ObjectID))
	}
	return
}

func (c *Command) InsertOrUpdate(document interface{}) (result CmdResult, err error) {
	model := mongo.NewUpdateOneModel().SetFilter(c.condition).SetUpsert(true).SetUpdate(document)
	res, err := c.collection.BulkWrite(context.Background(), []mongo.WriteModel{model})
	if err != nil {
		return result, err
	}
	if len(res.UpsertedIDs) > 0 {
		result.ID = res.UpsertedIDs[0].(primitive.ObjectID)
		result.Count = res.UpsertedCount
	} else {
		result.Count = res.ModifiedCount
	}
	return
}

func (c *Command) Replace(document interface{}) (result CmdResult, err error) {
	var res *mongo.UpdateResult
	res, err = c.collection.ReplaceOne(context.Background(), c.condition, document)
	if err != nil {
		return
	}
	result.Count = res.ModifiedCount
	return
}

// 對整個 document 修改
// filter :條件
// document: value 值
func (c *Command) Update(document ...interface{}) (result CmdResult, err error) {
	var data interface{}
	if c.updateModel.updateOperator != nil {
		data = c.updateModel.updateOperator
	}
	if len(document) > 0 {
		data = document[0]
	}
	if data == nil {
		err = errors.New("no update data")
		return
	}
	model := mongo.NewUpdateOneModel().SetFilter(c.condition).SetUpdate(data)
	res, err := c.collection.BulkWrite(context.Background(), []mongo.WriteModel{model})
	if err != nil {
		return result, err
	}
	result.Count = res.ModifiedCount
	return
}

func (c *Command) UpdateMany(document ...interface{}) (result CmdResult, err error) {
	var data interface{}
	if c.updateModel.updateOperator != nil {
		data = c.updateModel.updateOperator
	}
	if len(document) > 0 {
		data = document[0]
	}
	if data == nil {
		err = errors.New("no update data")
		return
	}
	model := mongo.NewUpdateManyModel().SetFilter(c.condition).SetUpdate(data)
	res, err := c.collection.BulkWrite(context.Background(), []mongo.WriteModel{model})
	if err != nil {
		return result, err
	}
	result.Count = res.ModifiedCount
	return
}

func (c *Command) Delete() (result CmdResult, err error) {
	var res *mongo.DeleteResult
	res, err = c.collection.DeleteMany(context.Background(), c.condition)
	if err != nil {
		return
	}
	result = CmdResult{Count: res.DeletedCount}
	return
}

func (c *Command) Count() (count int64, err error) {
	count, err = c.collection.CountDocuments(context.Background(), c.condition, c.PrepareCountOptions())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Command) PrepareCountOptions() *options.CountOptions {
	opts := options.Count()
	if c.limit != nil {
		opts.SetLimit(*c.limit)
	}

	if c.skip != nil {
		opts.SetSkip(*c.skip)
	}
	return opts
}

func (c *Command) Select(slicePtr interface{}) (page *Page, err error) {
	slicePtrV := reflect.ValueOf(slicePtr)
	if kind := slicePtrV.Kind(); kind != reflect.Ptr {
		err = fmt.Errorf("slicePtr type is not pointer,it is %s", kind.String())
		return
	}

	if slicePtrV.IsNil() {
		err = fmt.Errorf("slicePtr is nil")
		return
	}

	sliceV := slicePtrV.Elem()
	if kind := sliceV.Kind(); kind != reflect.Slice {
		err = fmt.Errorf("sliceVal type is not slice,it is %s", kind.String())
		return
	}

	var cur *mongo.Cursor
	cur, err = c.collection.Find(context.Background(), c.condition, c.prepareFindOptions())
	if err != nil {
		return
	}

	err = cur.All(context.Background(), slicePtr)
	if err != nil {
		return
	}
	return c.pageTotal()
}

func (c *Command) SelectRtnMap() (data []map[string]interface{}, err error) {
	var cur *mongo.Cursor
	cur, err = c.collection.Find(context.Background(), c.condition, c.prepareFindOptions())
	if err != nil {
		return nil, err
	}
	err = cur.All(context.Background(), &data)
	return data, nil
}

func (c *Command) prepareFindOptions() *options.FindOptions {
	opts := options.Find()
	if len(c.sort) > 0 {
		opts.SetSort(c.sort)
	}

	if len(c.projection) > 0 {
		opts.SetProjection(c.projection)
	}

	if c.limit != nil {
		opts.SetLimit(*c.limit)
	}

	if c.skip != nil {
		opts.SetSkip(*c.skip)
	}
	return opts
}

func (c *Command) SelectOne(data interface{}) (err error) {
	listType := reflect.TypeOf(data)
	if listType.Kind() == reflect.Ptr {

	} else {
		err = errors.New("unsupported destination, only pointer is supported on SelectOne")
		return
	}

	cur := c.collection.FindOne(context.Background(), c.condition, c.prepareFindOneOptions())
	if err = cur.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	if err = cur.Decode(data); err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}
	return
}

func (c *Command) prepareFindOneOptions() *options.FindOneOptions {
	opts := options.FindOne()
	if len(c.sort) > 0 {
		opts.SetSort(c.sort)
	}

	if len(c.projection) > 0 {
		opts.SetProjection(c.projection)
	}

	if c.skip != nil {
		opts.SetSkip(*c.skip)
	}
	return opts
}

func (c *Command) Aggregate(slicePtr interface{}) (page *Page, err error) {
	slicePtrV := reflect.ValueOf(slicePtr)
	if kind := slicePtrV.Kind(); kind != reflect.Ptr {
		err = fmt.Errorf("slicePtr type is not pointer,it is %s", kind.String())
		return
	}

	if slicePtrV.IsNil() {
		err = fmt.Errorf("slicePtr is nil")
		return
	}
	sliceV := slicePtrV.Elem()
	if kind := sliceV.Kind(); kind != reflect.Slice {
		err = fmt.Errorf("sliceVal type is not slice,it is %s", kind.String())
		return
	}

	var cur *mongo.Cursor
	cur, err = c.collection.Aggregate(context.Background(), c.pipeline, options.Aggregate())
	if err != nil {
		return
	}
	err = cur.All(context.Background(), slicePtr)
	if err != nil {
		return
	}
	return c.pageTotal()
}

func (c *Command) AggregateRtnMap() (data []map[string]interface{}, err error) {
	var cur *mongo.Cursor
	cur, err = c.collection.Aggregate(context.Background(), c.pipeline, options.Aggregate())
	if err != nil {
		return
	}

	err = cur.All(context.Background(), &data)
	return
}
