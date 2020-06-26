package bsontool

import "go.mongodb.org/mongo-driver/bson"

type M = bson.M

func Sum(m interface{}) M {
	return M{"$sum": m}
}

func Avg(m interface{}) M {
	return M{"$avg": m}
}

func Push(m interface{}) M {
	return M{"$push": m}
}

func Max(m interface{}) M {
	return M{"$max": m}
}

func Min(m interface{}) M {
	return M{"$min": m}
}
