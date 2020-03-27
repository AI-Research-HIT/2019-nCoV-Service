package service

import (
	"encoding/json"
	"errors"
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
	rs, err := db.FindLatestOverallData()
	if err != nil || rs.UpdateTime < now.Unix()*1000 {
		ewlog.Errorf("not find in db %v", err)

		all, err := cli.GetOverAll(1)
		if err != nil || len(all) == 0 {
			ewlog.Error(err)
		} else if len(all) > 0 {
			err = db.InsertLatestOverallData(all[0])
			if err != nil {
				ewlog.Error(err)
				resputil.WriteFailed(w, 400, "latest data not found")
				return
			}
			rs = all[0]
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

	// province, err := db.CalculateDataByDay(request.Province, request.City)

	// if len(province.Detail) == 0 {
	// 	ewlog.Warn("请求的城市没有初始数据")
	// 	resputil.WriteFailed(w, 101, "请求的城市没有初始数据")
	// 	return
	// }

	// DayOne := province.Detail[0]
	// startDate, err := time.Parse(util.DateLayout, DayOne.Date)
	// if err != nil {
	// 	ewlog.Error(err)
	// 	resputil.WriteFailed(w, 103, "internal error")
	// 	return
	// }
	startDate := time.Now()

	allResult := []protodef.MonteCarloSimulationResp{}

	calNum := 20

	if request.SimulateNum > 0 {
		calNum = request.SimulateNum
	}
	isCalSpread := true

	treatmentEffect := map[int]float64{
		0: 1.0,
	}

	var countTime int64 = 0

	errCount := 0

	for i := 0; i < calNum; i++ {
		if errCount > 3 {
			err = errors.New("该参数计算时间太长，疫情已经发展到超出30万人同时感染，取消计算, 请重新调整参数")
			ewlog.Error(err)
			resputil.WriteFailed(w, 103, err.Error())
			return
		}
		if i != 0 {
			isCalSpread = false
		}
		startT := time.Now()

		result, err := model.Simulate(request.InitNum, request.PredictDay, request.Mlist, request.BetaList,
			request.TreamentList, treatmentEffect, startDate, isCalSpread, request.IsQuarantineCloser,
			request.MedicalNum, true, request.NoSymptomProb)

		endT := time.Now()

		countTime += endT.UnixNano() - startT.UnixNano()
		if err != nil {
			ewlog.Error(err)
			errCount++
			//resputil.WriteFailed(w, 102, err.Error())
			continue
		}

		allResult = append(allResult, result)
		if countTime > int64(time.Second*8) {
			break
		}
	}

	ewlog.Infof("平局每次模拟时间%dms", countTime/int64(calNum)/int64(time.Millisecond))
	if len(allResult) == 0 {
		err = errors.New("该参数计算时间太长，疫情已经发展到超出30万人同时感染，取消计算, 请重新调整参数")
		ewlog.Error(err)
		resputil.WriteFailed(w, 103, err.Error())
		return
	}
	result := protodef.MonteCarloSimulationResp{
		SpreadTrack: allResult[0].SpreadTrack,
		Statistic:   []protodef.MonteCarloSimulationItem{},
	}

	calNum = len(allResult)

	if len(result.SpreadTrack.Nodes) > 500 {
		result.SpreadTrack = protodef.SpreadTrackResponse{}
	}

	for i := 0; i < len(allResult[0].Statistic)-1; i++ {
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
		InfectedNotQuarantineCount := 0
		CloserQuarantineCount := 0
		InfectedMin := 0
		InfectedMax := 0
		AttackedNotTreamenting := 0
		DeadProb := 0.0
		ConfirmingCount := 0
		NoSymptomCount := 0
		for j := 0; j < len(allResult); j++ {
			count := allResult[j].Statistic[i].InfectedCount
			if j == 0 {
				InfectedMin = count
				InfectedMax = count
			} else {
				if InfectedMin > count {
					InfectedMin = count
				}
				if InfectedMax < count {
					InfectedMax = count
				}
			}
			InfectedCount += count
			InfectedNew += allResult[j].Statistic[i].InfectedNew
			ConfirmCount += allResult[j].Statistic[i].ConfirmCount
			ConfirmNew += allResult[j].Statistic[i].ConfirmNew
			CureCount += allResult[j].Statistic[i].CureCount
			CureNew += allResult[j].Statistic[i].CureNew
			DeadCount += allResult[j].Statistic[i].DeadCount
			DeadNew += allResult[j].Statistic[i].DeadNew
			InfectingCount += allResult[j].Statistic[i].InfectingCount
			TreamentingCount += allResult[j].Statistic[i].TreamentingCount
			InfectedNotQuarantineCount += allResult[j].Statistic[i].InfectedNotQuarantineCount
			CloserQuarantineCount += allResult[j].Statistic[i].CloserQuarantineCount
			AttackedNotTreamenting += allResult[j].Statistic[i].AttackedNotTreamenting
			DeadProb += allResult[j].Statistic[i].DeadProb
			ConfirmingCount += allResult[j].Statistic[i].ConfirmingCount
			NoSymptomCount += allResult[j].Statistic[i].NoSymptomCount
		}

		count := InfectedCount / calNum
		min := count - (count-InfectedMin)/5
		max := (InfectedMax-count)/10 + count

		s := protodef.MonteCarloSimulationItem{
			InfectedCount:              count,
			ConfirmCount:               ConfirmCount / calNum,
			CureCount:                  CureCount / calNum,
			DeadCount:                  DeadCount / calNum,
			InfectedNew:                InfectedNew / calNum,
			ConfirmNew:                 ConfirmNew / calNum,
			CureNew:                    CureNew / calNum,
			DeadNew:                    DeadNew / calNum,
			InfectingCount:             InfectingCount / calNum,
			TreamentingCount:           TreamentingCount / calNum,
			InfectedNotQuarantineCount: InfectedNotQuarantineCount / calNum,
			CloserQuarantineCount:      CloserQuarantineCount / calNum,
			InfectedMin:                min,
			InfectedMax:                max - min,
			AttackedNotTreamenting:     AttackedNotTreamenting / calNum,
			DeadProb:                   DeadProb / float64(calNum),
			Date:                       allResult[0].Statistic[i].Date,
			ConfirmingCount:            ConfirmingCount / calNum,
			NoSymptomCount:             NoSymptomCount / calNum,
		}
		result.Statistic = append(result.Statistic, s)
	}

	// for i := range result.Statistic {
	// 	if i == len(province.Detail) {
	// 		break
	// 	}

	// 	result.Statistic[i].RealConfirmCount = province.Detail[i].TotalInfection
	// 	result.Statistic[i].RealConfirmNew = province.Detail[i].NewInfection
	// 	result.Statistic[i].RealCureCount = province.Detail[i].TotalCure
	// 	result.Statistic[i].RealCureNew = province.Detail[i].NewCure
	// 	result.Statistic[i].RealDeadCount = province.Detail[i].TotalDeath
	// 	result.Statistic[i].RealDeadNew = province.Detail[i].NewDeath
	// }

	resputil.WriteSuccessWithData(w, result)
}

func GetAllProvinceOrCountryName(w http.ResponseWriter, r *http.Request) {
	rs, err := db.FindAllProvinceOrCountry()
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 500, "internal error")
	}

	resputil.WriteSuccessWithData(w, rs)
}

func GetProvinceDailyData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	country := r.FormValue("region")

	if len(country) == 0 {
		resputil.WriteFailed(w, 100, "没有指定地区名")
		return
	}

	province, err := db.CalculateDataByDay(country, 0)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 101, "没有该地区数据")
		return
	}

	resputil.WriteSuccessWithData(w, province)
}

func CompareCountryTrending(w http.ResponseWriter, r *http.Request) {
	countrys := []string{
		"湖北省",
		"美国",
		"意大利",
		"伊朗",
		"英国",
		"德国",
		"日本",
		"韩国",
		"新加坡",
		"西班牙",
		"法国",
	}

	countryData := []protodef.ProvinceData{}

	for _, c := range countrys {
		data, err := db.CalculateDataByDay(c, 0)
		if err != nil {
			ewlog.Error(err)
			continue
		}
		countryData = append(countryData, data)
	}

	result := []protodef.ProvinceData{}
	max := 0
	for _, c := range countryData {
		p := protodef.ProvinceData{
			ProvinceName: c.ProvinceName,
			Detail:       []protodef.DetailData{},
		}
		for _, d := range c.Detail {
			if d.TotalInfection > 100 {
				p.Detail = append(p.Detail, d)
			}
		}

		if len(p.Detail) > max {
			max = len(p.Detail)
		}

		result = append(result, p)
	}

	resp := protodef.CompareCountryResp{
		Countrys: countrys,
		Data:     result,
		Max:      max,
	}

	resputil.WriteSuccessWithData(w, resp)
}

type MResultT struct {
	M      int
	Result []float64
}

func CompareWithDiffMHandler(w http.ResponseWriter, r *http.Request) {
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

	// province, err := db.CalculateDataByDay(request.Province, request.City)

	// if len(province.Detail) == 0 {
	// 	ewlog.Warn("请求的城市没有初始数据")
	// 	resputil.WriteFailed(w, 101, "请求的城市没有初始数据")
	// 	return
	// }

	// DayOne := province.Detail[0]
	// startDate, err := time.Parse(util.DateLayout, DayOne.Date)
	// if err != nil {
	// 	ewlog.Error(err)
	// 	resputil.WriteFailed(w, 103, "internal error")
	// 	return
	// }
	startDate := time.Now()

	calNum := 20

	if request.SimulateNum > 0 {
		calNum = request.SimulateNum
	}

	treatmentEffect := map[int]float64{
		0: 1.0,
	}

	var countTime int64 = 0

	errCount := 0

	mlist := []map[int]float64{}
	for i := 5; i <= 20; i++ {
		mlist = append(mlist, map[int]float64{0: float64(i)})
	}

	mresult := map[int]MResultT{}
	fresult := [][]float64{}

	initNumList := []int{1, 2, 4, 8, 16}
	//rsCh := make(chan MResultT, 15)

	for _, m := range mlist {
		tmpM := m
		tmpR := []float64{}
		//go func() {
		ewlog.Infof("Start: %v", tmpM)
		for _, num := range initNumList {

			allResult := []protodef.MonteCarloSimulationResp{}

			for i := 0; i < calNum; i++ {
				//ewlog.Infof("%v calNum: %d", tmpM, i)

				startT := time.Now()

				result, err := model.Simulate(num, request.PredictDay, tmpM, request.BetaList,
					request.TreamentList, treatmentEffect, startDate, false, request.IsQuarantineCloser,
					request.MedicalNum, true, request.NoSymptomProb)

				endT := time.Now()

				countTime += endT.UnixNano() - startT.UnixNano()
				if err != nil {
					//ewlog.Error(err)
					errCount++
					//resputil.WriteFailed(w, 102, err.Error())
					continue
				}

				allResult = append(allResult, result)
			}

			ewlog.Infof("%v, 平均每次模拟时间%dms", tmpM, countTime/int64(calNum)/int64(time.Millisecond))
			// if len(allResult) == 0 {
			// 	err = errors.New("该参数计算时间太长，疫情已经发展到超出30万人同时感染，取消计算, 请重新调整参数")
			// 	ewlog.Error(err)
			// 	resputil.WriteFailed(w, 103, err.Error())
			// 	return
			// }
			count := 0
			for _, s := range allResult {
				if s.Statistic[len(s.Statistic)-1].ConfirmNew == 0 {
					count++
				}
			}
			ewlog.Infof("%v 初始感染人数: %d 0百分比: %d", m, num, count)
			tmpR = append(tmpR, float64(count)/10.0)
		}
		mrs := MResultT{
			M:      int(tmpM[0]),
			Result: tmpR,
		}
		//rsCh <- mrs
		//}()
		mresult[mrs.M] = mrs
	}
	// count := 0
	// for i := range rsCh {
	// 	fmt.Printf("receiver get %+v\n", i)
	// 	mresult[i.M] = i
	// 	count++
	// 	if count >= 15 {
	// 		// signal recving finish
	// 		break
	// 	}
	// }

	for i := 5; i <= 20; i++ {
		re := mresult[i]
		fresult = append(fresult, re.Result)
	}

	// for i := range result.Statistic {
	// 	if i == len(province.Detail) {
	// 		break
	// 	}

	// 	result.Statistic[i].RealConfirmCount = province.Detail[i].TotalInfection
	// 	result.Statistic[i].RealConfirmNew = province.Detail[i].NewInfection
	// 	result.Statistic[i].RealCureCount = province.Detail[i].TotalCure
	// 	result.Statistic[i].RealCureNew = province.Detail[i].NewCure
	// 	result.Statistic[i].RealDeadCount = province.Detail[i].TotalDeath
	// 	result.Statistic[i].RealDeadNew = province.Detail[i].NewDeath
	// }

	resputil.WriteSuccessWithData(w, fresult)
}
