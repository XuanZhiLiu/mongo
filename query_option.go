package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (c *Command) Sort(column string, isDesc bool) *Command {
	var sortBy int32
	if isDesc {
		sortBy = -1
	} else {
		sortBy = 1
	}
	c.sort = append(c.sort, bsonx.Elem{
		column, bsonx.Int32(sortBy),
	})
	return c
}

func (c *Command) Projection(column string, isShow bool) *Command {
	var show int32
	if isShow {
		show = 1
	} else {
		show = 0
	}
	c.projection = append(c.projection, bson.E{
		Key: column, Value: show,
	})
	return c
}

func (c *Command) Projections(column []string, isShow bool) *Command {
	var show int32
	if isShow {
		show = 1
	} else {
		show = 0
	}
	for _, col := range column {
		c.projection = append(c.projection, bson.E{
			Key: col, Value: show,
		})
	}
	return c
}

func (c *Command) ProjectionsArrayMatch(parent, column string, value interface{}) *Command {
	c.projection = append(c.projection, bson.E{
		Key: parent, Value: bson.M{"$elemMatch": bson.D{{column, value}}},
	})
	return c
}

func (c *Command) ProjectionsArrayIndex(column string, skip, limit int) *Command {
	c.projection = append(c.projection, bson.E{
		Key: column, Value: bson.D{{"$slice", []int{skip, limit}}},
	})
	return c
}

func (c *Command) Limit(limit int64) *Command {
	c.limit = &limit
	return c
}

func (c *Command) Skip(skip int64) *Command {
	c.skip = &skip
	return c
}

func (c *Command) Page(index, size int64) *Command {
	if index <= 0 {
		index = 1
	}
	if size <= 0 {
		size = 10
	}
	skip := (index - 1) * size
	c.skip = &skip
	c.limit = &size
	return c
}

func (c *Command) Pipeline(pipeline Pipeline) *Command {
	c.pipeline = pipeline
	return c
}
