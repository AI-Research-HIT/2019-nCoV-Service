package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

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
		resputil.WriteFailed(w, 2, err.Error())
		return
	}

	defer r.Body.Close()

	request := &protodef.PredictionRequest{}
	err = json.Unmarshal(data, request)

	resp := model.Prediction(*request)

	resputil.WriteSuccessWithData(w, resp)
}

func MDataHandler(w http.ResponseWriter, r *http.Request) {
	resputil.WriteSuccessWithData(w, model.TempMData)
}

func LatestDataHandler(w http.ResponseWriter, r *http.Request) {
	data := protodef.LatestData{
		NowInfection:   54645,
		TotalCure:      18687,
		TotalDeath:     2239,
		NowHeavy:       11633,
		NowSusp:        5206,
		TotalInfection: 75571,
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
