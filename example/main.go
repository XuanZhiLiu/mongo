package main

import (
	"fmt"
	"time"

	"github.com/XuanZhiLiu/mongo"
	"github.com/XuanZhiLiu/mongo/bsontool"
)

type mongoDB struct {
	User     string
	Password string
	Host     string
	Dbname   string
}

var dbConf = &mongoDB{
	User:     "",
	Password: "",
	Host:     "10.200.252.114:27017",
	Dbname:   "test",
}

type PersonModel struct {
	mongo.CollectionBase `bson:",inline"`
	Name                 string        `json:"name" bson:"name"`
	Birth                time.Time     `json:"birth" bson:"birth"`
	Age                  int           `json:"age" bson:"age"`
	Tall                 int           `json:"tall" bson:"tall"`
	FavorNum             []int         `json:"favorNum" bson:"favorNum"`
	Parents              []PersonModel `json:"parents" bson:"parents"`
}

var testPeopleData = []PersonModel{
	{
		Name:     "Kevin",
		Birth:    time.Date(1991, 2, 10, 15, 25, 0, 0, time.UTC),
		Age:      29,
		Tall:     170,
		FavorNum: []int{1, 2, 3, 5},
		Parents: []PersonModel{
			{
				Name:  "Peter",
				Birth: time.Date(1970, 5, 20, 4, 14, 0, 0, time.UTC),
				Age:   55,
			},
			{
				Name:  "Tina",
				Birth: time.Date(1972, 4, 15, 6, 14, 0, 0, time.UTC),
				Age:   53,
			},
		},
	},
	{
		Name:     "George",
		Birth:    time.Date(1994, 7, 25, 18, 0, 0, 0, time.UTC),
		Age:      25,
		Tall:     177,
		FavorNum: []int{2, 3, 4},
		Parents: []PersonModel{
			{
				Name:  "Jack",
				Birth: time.Date(1966, 8, 20, 4, 14, 0, 0, time.UTC),
				Age:   58,
			},
			{
				Name:  "Emy",
				Birth: time.Date(1970, 4, 15, 6, 14, 0, 0, time.UTC),
				Age:   55,
			},
		},
	},
}

func main() {
	ConnectMongoDb()
	Test()
	fmt.Println(mongo.ListCollections(nil))
	fmt.Println(mongo.ListCollections(bsontool.D{{"name", "testIndex"}}))

	InsertIndex()

	fmt.Println(mongo.NewCommand(nil, "person").DropCollection())
	InsertManyData()

	fmt.Println("ExampleInsertOrUpdate")
	ExampleInsertOrUpdate()
	fmt.Println("")

	fmt.Println("ExampleGetEqual")
	ExampleGetEqual()
	fmt.Println("")

	fmt.Println("ExampleUpdate")
	ExampleUpdate()
	fmt.Println("")

	fmt.Println("ExampleGetEquals")
	ExampleGetEquals()
	fmt.Println("")

	fmt.Println("ExampleUpdateField")
	ExampleUpdateField()
	fmt.Println("")

	fmt.Println("ExampleGetIn")
	ExampleGetIn()
	fmt.Println("")

	fmt.Println("ExampleGetNotIn")
	ExampleGetNotIn()
	fmt.Println("")

	fmt.Println("ExampleGetAll")
	ExampleGetAll()
	fmt.Println("")

	fmt.Println("ExampleUpdateFields")
	ExampleUpdateFields()
	fmt.Println("")

	fmt.Println("ExampleGetSizeOfArray")
	ExampleGetSizeOfArray()
	fmt.Println("")

	fmt.Println("ExampleGetRegex")
	ExampleGetRegex()
	fmt.Println("")

	fmt.Println("ExampleGetNotEqual")
	ExampleGetNotEqual()
	fmt.Println("")

	fmt.Println("ExampleGetExists")
	ExampleGetExists()
	fmt.Println("")

	fmt.Println("ExampleGetBetween")
	ExampleGetBetween()
	fmt.Println("")

	fmt.Println("ExampleGetOne")
	ExampleGetOne()
	fmt.Println("")

	ExampleAggregate()

	fmt.Println("ExamplePull")
	ExamplePull()
}

// mongoDB 連線
func ConnectMongoDb() {
	mongo.Init(dbConf.Host, dbConf.User, dbConf.Password, dbConf.Dbname)
}

func InsertIndex() {
	mongo.NewCommand(nil, "testIndex").DropCollection()

	indexArray := []mongo.InsertIndexModel{}

	index := mongo.InsertIndexModel{}
	index.Keys = bsontool.D{
		{Key: "code", Value: 1},
		{Key: "period", Value: -1},
		{Key: "drawTime", Value: -1},
		{Key: "status", Value: 1},
	}
	index.Options = mongo.NewIndexOption().SetName("testIndexIndex1").SetUnique(true).
		SetExpireAfterSeconds(120).SetBackground(true).SetSparse(true)
	indexArray = append(indexArray, index)

	index = mongo.InsertIndexModel{}
	index.Keys = bsontool.D{
		{Key: "code", Value: 1},
		{Key: "period", Value: -1},
		{Key: "drawTime", Value: -1},
	}
	index.Options = mongo.NewIndexOption().SetName("testIndexIndex2")
	indexArray = append(indexArray, index)

	index = mongo.InsertIndexModel{}
	index.Keys = bsontool.D{
		{Key: "code", Value: 1},
		{Key: "period", Value: -1},
	}
	index.Options = mongo.NewIndexOption().SetName("testIndexIndex3")
	indexArray = append(indexArray, index)

	fmt.Println(mongo.NewCommand(nil, "testIndex").CreateIndexes(indexArray...))
	fmt.Println(mongo.NewCommand(nil, "testIndex").IndexList())
	fmt.Println(mongo.NewCommand(nil, "testIndex").DropIndex("testIndexIndex1"))
	fmt.Println(mongo.NewCommand(nil, "testIndex").IndexList())
	fmt.Println(mongo.NewCommand(nil, "testIndex").DropAllIndex())
	fmt.Println(mongo.NewCommand(nil, "testIndex").IndexList())
	fmt.Println(mongo.NewCommand(nil, "testIndex").CreateIndexes(indexArray...))

}

func InsertManyData() {
	mongo.NewCommand(nil, "person").Delete()
	for _, people := range testPeopleData {
		result, err := mongo.NewCommand(nil, "person").Insert(people)
		if err != nil {
			fmt.Println("SQL", err.Error())
			return
		}
		if result.ID.IsZero() {
			fmt.Println("no row inserted")
		}
	}
	return
}

func ExampleUpdate() {
	data := testPeopleData[1]
	data.Age = 24
	result, err := mongo.NewCommand(nil, "person").Equal("name", true, "George").Update(data)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	fmt.Println(result)
}

func ExampleUpdateField() {
	result, err := mongo.NewCommand(nil, "person").Equal("name", true, "Kevin").Set("age", 30).Update()
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	fmt.Println(result)
}

func ExampleUpdateFields() {

	result, err := mongo.NewCommand(nil, "person").Equal("name", true, "Kevin").Set("tall", 171).Set("birth", time.Date(1991, 2, 10, 16, 25, 0, 0, time.UTC)).Update()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func ExampleInsertOrUpdate() {
	data := testPeopleData[0]
	data.Name = "Kay"
	data.Tall = 172
	result, err := mongo.NewCommand(nil, "person").Equal("name", true, "Kay").InsertOrUpdate(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func ExampleGetEqual() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").Equal("name", true, "George").Equal("age", true, 25).Select(&result)
	// result, err := mongo.NewCommand(nil, "person").Equal("name", true, "George").Equal("age", true, 25).SelectRtnMap()
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 1 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetEquals() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").Equals(true, bsontool.KeyValue{"name", "Kevin"}, bsontool.KeyValue{"age",
		29}).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 1 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetIn() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").In("age", true, 25, 30).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 1 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetNotIn() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").NotIn("name", true, "Kevin", "Georgee").Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 2 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetAll() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").All("favorNum", true, 2, 3, 4).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 1 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetSizeOfArray() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").SizeOfArray("favorNum", true, 4).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 2 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetRegex() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").Regex("name", true, "^K.*", "").Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 2 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetNotEqual() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").NotEqual("age", true, 29).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 2 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetExists() {
	var result []PersonModel
	_, err := mongo.NewCommand(nil, "person").Exists("age", true, true).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 3 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetBetween() {
	var result []PersonModel
	var gt *int
	_, err := mongo.NewCommand(nil, "person").Between("age", true, 29, gt, nil, nil).Select(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	if len(result) == 2 {
		fmt.Println(result)
	} else {
		fmt.Println("fail")
	}
}

func ExampleGetOne() {
	var result *PersonModel
	err := mongo.NewCommand(nil, "person").Equal("name", true, "George").SelectOne(&result)
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	fmt.Println(result)
}

func ExampleAggregate() {
	var pipelines []mongo.Pipeline
	pipelines = append(pipelines, mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("name", true, "George")))

	pipelines = append(pipelines, mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("name", true, "George")).Group("$name", bsontool.M{"count": bsontool.M{"$sum": 1}}))

	pipelines = append(pipelines, mongo.NewPipeline().Group(bsontool.M{"$dateToString": bsontool.M{"format": "%Y-%m-%d", "date": "$birth"}}, bsontool.M{"count": bsontool.M{"$sum": 1}}))

	pipelines = append(pipelines, mongo.NewPipeline().Group(bsontool.M{"$dateToString": bsontool.M{"format": "%Y-%m-%d", "date": "$birth"}}, bsontool.M{"count": bsontool.Sum(bsontool.M{"$multiply": []string{"$age", "$tall"}})}))

	pipelines = append(pipelines, mongo.NewPipeline().Group(bsontool.M{"$dateToString": bsontool.M{"format": "%Y-%m-%d", "date": "$birth"}}, bsontool.M{"count": bsontool.Avg("$age")}))

	pipelines = append(pipelines, mongo.NewPipeline().Group(bsontool.M{"$dateToString": bsontool.M{"format": "%Y-%m-%d", "date": "$birth"}}, bsontool.M{"count": bsontool.Push("$age")}))

	pipelines = append(pipelines, mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("name", true, "George")).Project(bsontool.M{"age": 1}))

	pipelines = append(pipelines, mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("name", true, "George")).Unwind("$favorNum"))

	pipelines = append(pipelines, mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("name", true, "George")).Lookup("personLookup", "name", "name", "lookup_doc", nil, nil))

	pipelines = append(pipelines, mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("name", true, "George")).Count("name"))

	for i, pipeline := range pipelines {
		if i == 0 {
			var data []*PersonModel
			_, err := mongo.NewCommand(nil, "person").Pipeline(pipeline).Aggregate(&data)
			if err != nil {
				fmt.Println("SQL", err.Error())
				return
			}
			fmt.Println("index: ", i)
			fmt.Println(data)
			fmt.Println("2")
		} else {
			data, err := mongo.NewCommand(nil, "person").Pipeline(pipeline).AggregateRtnMap()
			if err != nil {
				fmt.Println("SQL", err.Error())
				return
			}
			fmt.Println("index: ", i)
			fmt.Println(data)
			fmt.Println("2")
		}
	}
}

func ExamplePull() {
	mongo.NewCommand(nil, "testArray").Pull("interests", "abc").Update()
}

func Test() {
	cmd := mongo.NewCommand(nil, "testUpdate").Set("abc", "123").Set("cde", "sre").Inc("12243", 1)
	_, err := cmd.Update()
	if err != nil {
		fmt.Println("SQL", err.Error())
	}

	pipe := mongo.NewPipeline().Match(bsontool.NewBsonD().Equal("hbsq_id", true, 1))
	// .Group("$name", bsontool.M{"count": bsontool.M{"$sum": 1}})
	data, err := mongo.NewCommand(nil, "testEdward").Pipeline(pipe).AggregateRtnMap()
	if err != nil {
		fmt.Println("SQL", err.Error())
		return
	}
	_ = data
}

func Abc(data bsontool.M) {

}
