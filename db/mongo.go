package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/cli"
	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/AI-Research-HIT/2019-nCoV-Service/util"
	"github.com/ender-wan/ewlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName          = "2019-nCoV"
	epidCol         = "epid"
	qianxiCol       = "baiduQianxi"
	overallCol      = "overall"
	provinceDataCol = "provinceData"
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
		return
	}

	fmt.Println("MongoDB connected.")

	go func() {
		util.Recover()

		err = InitOverallData()
		if err != nil {
			ewlog.Error(err)
		}

		err = InitAllProvinceData()
		if err != nil {
			ewlog.Error(err)
		}
	}()

	go func() {
		util.Recover()
		ticker := time.NewTicker(time.Hour * 1)

		for {
			select {
			case <-ticker.C:

			}
			ewlog.Info("fetch latest data")
			err = fetchLatestProvinceData()
			if err != nil {
				ewlog.Error(err)
			}
		}
	}()

}

func FindProvinceData(pid string) (provinceData protodef.ProvinceData, err error) {
	ctx := context.Background()
	result := dbClient.Database(dbName).Collection(epidCol).FindOne(ctx, bson.M{"pid": pid})
	err = result.Decode(&provinceData)

	return
}

func InsertLatestOverallData(overall cli.OverAllT) (err error) {
	ctx := context.Background()

	_, err = dbClient.Database(dbName).Collection(overallCol).InsertOne(ctx, overall)

	return err
}

func FindLatestOverallData() (result cli.OverAllT, err error) {
	ctx := context.Background()

	rs := dbClient.Database(dbName).Collection(overallCol).FindOne(ctx, bson.M{})

	err = rs.Decode(&result)

	return
}

func InsertProvinceData(province cli.ProvinceT) (err error) {
	ctx := context.Background()

	_, err = dbClient.Database(dbName).Collection(provinceDataCol).InsertOne(ctx, province)

	return err
}

func FindAllProvinceData(provinceName string) (data []cli.ProvinceT, err error) {
	ctx := context.Background()

	result, err := dbClient.Database(dbName).Collection(provinceDataCol).Find(ctx, bson.M{"provinceName": provinceName})
	if err != nil {
		ewlog.Error(err)
		return
	}

	for result.Next(ctx) {
		var province cli.ProvinceT
		if err = result.Decode(&province); err == nil {
			data = append(data, province)
		} else {
			ewlog.Error(err)
		}
	}

	return
}

func InitOverallData() (err error) {
	ctx := context.Background()
	count, err := dbClient.Database(dbName).Collection(overallCol).CountDocuments(ctx, bson.M{})
	//ewlog.Info(count)

	if count == 0 {
		rs, err := cli.GetOverAll(0)
		if err != nil {
			ewlog.Error(err)
			return err
		}
		ewlog.Info("fetch number: ", len(rs))
		for _, r := range rs {
			err = InsertLatestOverallData(r)
			if err != nil {
				ewlog.Error(err)
			}
		}
	}

	return
}

func InitAllProvinceData() (err error) {
	ctx := context.Background()

	rs, err := cli.GetProvinceNames()
	if err != nil {
		ewlog.Error(err)
		return err
	}
	ewlog.Info("fetch number: ", len(rs))
	for _, r := range rs {
		count, err := dbClient.Database(dbName).Collection(provinceDataCol).CountDocuments(ctx, bson.M{"provinceName": r})
		ewlog.Infof("%s: %d", r, count)
		if err != nil {
			ewlog.Error(err)
			continue
		}
		if count == 0 {
			all, err := cli.RetryGetAllProvinceData(0, r, 3)
			if err != nil {
				ewlog.Info("failed: ", r)
				ewlog.Error(err)
				continue
			}
			for _, p := range all {
				err = InsertProvinceData(p)
				if err != nil {
					ewlog.Error(err)
				}
			}
		}

	}

	return
}

func FindAllProvinceOrCountry() ([]string, error) {
	ctx := context.Background()

	rs, err := dbClient.Database(dbName).Collection(provinceDataCol).Distinct(ctx, "provinceName", bson.M{})
	if err != nil {
		ewlog.Error(err)
		return nil, err
	}

	all := []string{}
	for _, s := range rs {
		str := s.(string)
		all = append(all, str)
	}

	return all, err
}

func fetchLatestProvinceData() (err error) {
	ctx := context.Background()

	rs, err := cli.GetProvinceNames()
	if err != nil {
		ewlog.Error(err)
		return err
	}
	ewlog.Info("fetch number: ", len(rs))
	for _, r := range rs {
		all, err := cli.RetryGetAllProvinceData(1, r, 3)
		if err != nil {
			ewlog.Info("failed: ", r)
			ewlog.Error(err)
			continue
		}
		if len(all) > 0 {
			d := all[0]
			count, err := dbClient.Database(dbName).Collection(provinceDataCol).CountDocuments(ctx, bson.M{"provinceName": r, "updateTime": d.UpdateTime})
			ewlog.Infof("%s: %d", r, count)
			if err != nil {
				ewlog.Error(err)
				continue
			}

			if count == 0 {
				err = InsertProvinceData(all[0])
				if err != nil {
					ewlog.Error(err)
				}
			}
		}

		time.Sleep(time.Second * 3)
	}

	return
}
