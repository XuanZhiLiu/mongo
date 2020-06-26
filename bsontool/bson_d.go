package bsontool

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type D bson.D

func NewBsonD() D {
	return D{}
}

type KeyValue struct {
	Key   string
	Value interface{}
}

// enable:由client自行判斷是否需要帶此條件
func (c D) CustomConditions(cond D) D {
	c = cond
	return c
}

func (c D) Equals(enable bool, vals ...KeyValue) D {
	if !enable {
		return c
	}
	for _, val := range vals {
		c = append(c, bson.E{val.Key, val.Value})
	}
	return c
}

func (c D) Equal(col string, enable bool, val interface{}) D {
	if isZeroOfUnderlyingType(val) || !enable {
		return c
	}
	c = append(c, bson.E{col, val})
	return c
}

func (c D) NotEqual(col string, enable bool, val interface{}) D {
	if isZeroOfUnderlyingType(val) || !enable {
		return c
	}
	c = append(c, bson.E{col, bson.D{{"$ne", val}}})
	return c
}

func (c D) Exists(col string, enable, exists bool) D {
	if !enable {
		return c
	}
	c = append(c, bson.E{col, bson.D{{"$exists", exists}}})
	return c
}

func (c D) In(col string, enable bool, vals ...interface{}) D {
	if len(vals) == 0 || !enable {
		return c
	}
	c = append(c, bson.E{col, bson.D{{"$in", vals}}})
	return c
}

func (c D) NotIn(col string, enable bool, vals ...interface{}) D {
	if len(vals) == 0 || !enable {
		return c
	}
	c = append(c, bson.E{col, bson.D{{"$nin", vals}}})
	return c
}

func (c D) All(col string, enable bool, vals ...interface{}) D {
	if len(vals) == 0 || !enable {
		return c
	}
	c = append(c, bson.E{col, bson.D{{"$all", vals}}})
	return c
}

func (c D) SizeOfArray(col string, enable bool, size int) D {
	if !enable {
		return c
	}
	c = append(c, bson.E{col, bson.D{{"$size", size}}})
	return c
}

func (c D) Regex(col string, enable bool, pattern, options string) D {
	if !enable {
		return c
	}
	c = append(c, bson.E{col, primitive.Regex{Pattern: pattern, Options: options}})
	return c
}

func (c D) Between(col string, enable bool, gte, gt, lte, lt interface{}) D {
	if !enable {
		return c
	}
	conditions := make(bson.M)
	if !isZeroOfUnderlyingType(gte) {
		conditions["$gte"] = gte
	}

	if !isZeroOfUnderlyingType(gt) {
		conditions["$gt"] = gt
	}

	if !isZeroOfUnderlyingType(lte) {
		conditions["$lte"] = lte
	}

	if !isZeroOfUnderlyingType(lt) {
		conditions["$lt"] = lt
	}
	if len(conditions) == 0 {
		return c
	}
	c = append(c, bson.E{col, conditions})
	return c
}

func (c D) GreaterThanOrEqual(col string, enable bool, val interface{}) D {
	return c.Between(col, enable, val, nil, nil, nil)
}

func (c D) GreaterThan(col string, enable bool, val interface{}) D {
	return c.Between(col, enable, nil, val, nil, nil)
}

func (c D) LessThanOrEqual(col string, enable bool, val interface{}) D {
	return c.Between(col, enable, nil, nil, val, nil)
}

func (c D) LessThan(col string, enable bool, val interface{}) D {
	return c.Between(col, enable, nil, nil, nil, val)
}

func (c D) Or(conds ...interface{}) D {
	c = append(c, bson.E{"$or", conds})
	return c
}

func (c D) And(conds ...interface{}) D {
	c = append(c, bson.E{"$and", conds})
	return c
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
