package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/cli"

	"github.com/AI-Research-HIT/2019-nCoV-Service/db"

	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/AI-Research-HIT/2019-nCoV-Service/util"
	"github.com/ender-wan/ewlog"

	"github.com/AI-Research-HIT/2019-nCoV-Service/model"
	"github.com/AI-Research-HIT/2019-nCoV-Service/resputil"
)

func ModelCalculateHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 100, err.Error())
		return
	}

	defer r.Body.Close()

	request := &protodef.PredictionRequest{}
	err = json.Unmarshal(data, request)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 101, "无效的请求数据")
		return
	}
	ewlog.Infof("%+v", request)

	resp, err := model.Prediction(*request)
	if err != nil {
		resputil.WriteFailed(w, 102, "该区域数据太少，无法计算")
		return
	}

	resputil.WriteSuccessWithData(w, resp)
}

func MDataHandler(w http.ResponseWriter, r *http.Request) {
	resputil.WriteSuccessWithData(w, model.TempMData)
}

func LatestDataHandler(w http.ResponseWriter, r *http.Request) {
	now := util.TodayStartTime()
	rs, err := db.FindLatestOverallData(now.Unix())
	if err != nil {
		ewlog.Error(err)

		all, err := cli.GetOverAll(1)
		if err != nil {
			ewlog.Error(err)
			resputil.WriteFailed(w, 400, "latest data not found")
			return
		}
		if len(all) > 0 {
			err = db.InsertLatestOverallData(all[0])
			if err != nil {
				ewlog.Error(err)
				resputil.WriteFailed(w, 400, "latest data not found")
				return
			}
			rs = all[0]
		} else {
			ewlog.Errorf("overall api call failed")
			resputil.WriteFailed(w, 400, "latest data not found")
			return
		}

	}
	data := protodef.LatestData{
		NowInfection:   rs.CurrentConfirmedCount,
		TotalCure:      rs.CuredCount,
		TotalDeath:     rs.DeadCount,
		NowHeavy:       rs.SeriousCount,
		NowSusp:        rs.SuspectedCount,
		TotalInfection: rs.ConfirmedCount,
	}
	resputil.WriteSuccessWithData(w, data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := protodef.UserResponse{
		ID:       1,
		UserName: "admin",
		Password: "111",
		Avatar:   "http://47.75.202.128:8080/static/assets/user.png",
		Name:     "admin",
	}

	resputil.WriteSuccessWithData(w, user)
}

type MT struct {
	M    float64 `json:"m"`
	Date string  `json:"date"`
}

func MlistHandler(w http.ResponseWriter, r *http.Request) {
	mlist := []float64{16.0 / 16.0, 16.0 / 16.0, 13.0 / 16.0, 10.0 / 16.0, 8.0 / 16.0, 7.0 / 16.0, 6.0 / 16.0, 5.0 / 16.0, 5.0 / 16.0, 4.0 / 16.0, 4.0 / 16.0, 5.0 / 16.0, 3.0 / 16.0, 5.0 / 16.0, 4.0 / 16.0, 3.0 / 16.0, 2.0 / 16.0, 2.0 / 16.0, 2.0 / 16.0, 2.0 / 16.0, 2.0 / 16.0, 2.0 / 16}
	datelist := []string{"1-24", "1-25", "1-26", "1-27", "1-28", "1-29", "1-30", "1-31", "2-01", "2-02", "2-03", "2-04", "2-05", "2-06", "2-07", "2-08", "2-09", "2-10", "2-11", "2-12", "2-13", "2-14"}

	mtlist := []MT{}
	for i, _ := range mlist {
		mt := MT{
			M:    mlist[i],
			Date: datelist[i],
		}
		mtlist = append(mtlist, mt)
	}

	resputil.WriteSuccessWithData(w, mtlist)
}

type ChinaMobileT struct {
	Date string `json:"date"`
	Num  int    `json:"num"`
}

func ChinaMobileDataHandler(w http.ResponseWriter, r *http.Request) {
	data := []ChinaMobileT{}
	err := json.Unmarshal([]byte(model.ChinaMobileData), &data)
	if err != nil {
		ewlog.Error(err)
		return
	}
	resputil.WriteSuccessWithData(w, data)
}

func BaiduInCityHandler(w http.ResponseWriter, r *http.Request) {
	data := []protodef.BaiduCityT{}
	err := json.Unmarshal([]byte(model.BaiduIncity), &data)
	if err != nil {
		ewlog.Error(err)
		return
	}

	sort.Sort(util.BaiduCitySlice(data))

	resputil.WriteSuccessWithData(w, data)
}

func AllProvinceDataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 2, err.Error())
		return
	}

	defer r.Body.Close()

	request := &protodef.AllProviceDataRequest{}
	err = json.Unmarshal(data, request)

	resp, err := db.CalculateDataByDay(request.ProvinceName, request.CityName)
	if err != nil {
		resputil.WriteFailed(w, 404, request.ProvinceName+" data not found")
		return
	}

	resputil.WriteSuccessWithData(w, resp)
}

func MonteCarloSimulationHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 100, err.Error())
		return
	}

	defer r.Body.Close()

	request := &protodef.MonteCarloSimulationRequest{}
	err = json.Unmarshal(data, request)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 101, "无效的请求数据")
		return
	}
	ewlog.Infof("%+v", request)

	province, err := db.CalculateDataByDay(request.Province, request.City)

	if len(province.Detail) == 0 {
		ewlog.Warn("请求的城市没有初始数据")
		resputil.WriteFailed(w, 101, "请求的城市没有初始数据")
		return
	}

	DayOne := province.Detail[0]
	startDate, err := time.Parse(util.DateLayout, DayOne.Date)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 103, "internal error")
		return
	}

	allResult := []protodef.MonteCarloSimulationResp{}

	calNum := 20

	if request.SimulateNum > 0 {
		calNum = request.SimulateNum
	}

	for i := 0; i < calNum; i++ {
		result, err := model.Simulate(DayOne.TotalInfection, request.PredictDay, request.Mlist, request.BetaList, request.TreamentList, startDate)
		if err != nil {
			ewlog.Error(err)
			//resputil.WriteFailed(w, 102, err.Error())
			continue
		}
		allResult = append(allResult, result)
	}

	result := protodef.MonteCarloSimulationResp{
		SpreadTrack: allResult[0].SpreadTrack,
		Statistic:   []protodef.MonteCarloSimulationItem{},
	}

	for i := 0; i < len(allResult[0].Statistic); i++ {
		InfectedCount := 0
		InfectedNew := 0
		ConfirmCount := 0
		ConfirmNew := 0
		CureCount := 0
		CureNew := 0
		DeadCount := 0
		DeadNew := 0
		InfectingCount := 0
		TreamentingCount := 0
		for j := 0; j < len(allResult); j++ {
			InfectedCount += allResult[j].Statistic[i].InfectedCount
			InfectedNew += allResult[j].Statistic[i].InfectedNew
			ConfirmCount += allResult[j].Statistic[i].ConfirmCount
			ConfirmNew += allResult[j].Statistic[i].ConfirmNew
			CureCount += allResult[j].Statistic[i].CureCount
			CureNew += allResult[j].Statistic[i].CureNew
			DeadCount += allResult[j].Statistic[i].DeadCount
			DeadNew += allResult[j].Statistic[i].DeadNew
			InfectingCount += allResult[j].Statistic[i].InfectingCount
			TreamentingCount += allResult[j].Statistic[i].TreamentingCount
		}

		s := protodef.MonteCarloSimulationItem{
			InfectedCount:    InfectedCount / calNum,
			ConfirmCount:     ConfirmCount / calNum,
			CureCount:        CureCount / calNum,
			DeadCount:        DeadCount / calNum,
			InfectedNew:      InfectedNew / calNum,
			ConfirmNew:       ConfirmNew / calNum,
			CureNew:          CureNew / calNum,
			DeadNew:          DeadNew / calNum,
			InfectingCount:   InfectingCount / calNum,
			TreamentingCount: TreamentingCount / calNum,
			Date:             allResult[0].Statistic[i].Date,
		}
		result.Statistic = append(result.Statistic, s)

	}

	for i := range result.Statistic {
		if i == len(province.Detail) {
			break
		}

		result.Statistic[i].RealConfirmCount = province.Detail[i].TotalInfection
		result.Statistic[i].RealConfirmNew = province.Detail[i].NewInfection
		result.Statistic[i].RealCureCount = province.Detail[i].TotalCure
		result.Statistic[i].RealCureNew = province.Detail[i].NewCure
		result.Statistic[i].RealDeadCount = province.Detail[i].TotalDeath
		result.Statistic[i].RealDeadNew = province.Detail[i].NewDeath
	}

	resputil.WriteSuccessWithData(w, result)
}
