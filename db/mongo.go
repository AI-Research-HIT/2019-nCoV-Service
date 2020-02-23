package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/ender-wan/ewlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName    = "2019-nCoV"
	epidCol   = "epid"
	qianxiCol = "baiduQianxi"
)

var dbClient *mongo.Client

var HistorynCoVData = make(map[string]interface{})

func ConnectToMongo() {
	// jsonData, err := ioutil.ReadFile("data/2019-nCoV-data.json")
	// if err != nil {
	// 	ewlog.Error(err)
	// 	return
	// }
	// err = json.Unmarshal(jsonData, &HistorynCoVData)
	// if err != nil {
	// 	ewlog.Error(err)
	// 	return
	// }
	var err error
	var ctx context.Context
	var cancel context.CancelFunc

	dbClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		ewlog.Error(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		ewlog.Error(err)
	}
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		ewlog.Error(err)
	}

	err = dbClient.Ping(context.TODO(), nil)
	if err != nil {
		ewlog.Error(err)
	} else {
		fmt.Println("MongoDB connected.")
	}
}

func FindProvinceData(pid string) (provinceData protodef.ProvinceData, err error) {
	ctx := context.Background()
	result := dbClient.Database(dbName).Collection(epidCol).FindOne(ctx, bson.M{"pid": pid})
	err = result.Decode(&provinceData)

	return
}
