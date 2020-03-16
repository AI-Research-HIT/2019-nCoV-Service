package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ender-wan/ewlog"
)

type OverAllResultT struct {
	Results []OverAllT `json:"results" bson:"results"`
	Success bool       `json:"success" bson:"success"`
}

type OverAllT struct {
	CurrentConfirmedCount int    `json:"currentConfirmedCount" bson:"currentConfirmedCount"`
	ConfirmedCount        int    `json:"confirmedCount" bson:"confirmedCount"`
	SuspectedCount        int    `json:"suspectedCount" bson:"suspectedCount"`
	CuredCount            int    `json:"curedCount" bson:"curedCount"`
	DeadCount             int    `json:"deadCount" bson:"deadCount"`
	SeriousCount          int    `json:"seriousCount" bson:"seriousCount"`
	CurrentConfirmedIncr  int    `json:"currentConfirmedIncr" bson:"currentConfirmedIncr"`
	ConfirmedIncr         int    `json:"confirmedIncr" bson:"confirmedIncr"`
	SuspectedIncr         int    `json:"suspectedIncr" bson:"suspectedIncr"`
	CuredIncr             int    `json:"curedIncr" bson:"curedIncr"`
	DeadIncr              int    `json:"deadIncr" bson:"deadIncr"`
	SeriousIncr           int    `json:"seriousIncr" bson:"seriousIncr"`
	GeneralRemark         string `json:"generalRemark" bson:"generalRemark"`
	UpdateTime            int64  `json:"updateTime" bson:"updateTime"`
}

func HttpGet(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		ewlog.Error(err)
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		ewlog.Error(err)
	}

	return
}

// dtypt: 1 最新，0 所有
func GetOverAll(dtype int) (overall []OverAllT, err error) {
	url := fmt.Sprintf("https://lab.isaaclin.cn/nCoV/api/overall?latest=%d", dtype)

	body, err := HttpGet(url)

	if err != nil {
		ewlog.Error(err)
		return
	}

	result := OverAllResultT{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		ewlog.Error(err)
	}

	return result.Results, err
}

type ProvinceT struct {
	ProvinceName          string  `json:"provinceName" bson:"provinceName"`
	ProvinceShortName     string  `json:"provinceShortName" bson:"provinceShortName"`
	CurrentConfirmedCount int     `json:"currentConfirmedCount" bson:"currentConfirmedCount"`
	ConfirmedCount        int     `json:"confirmedCount" bson:"confirmedCount"`
	SuspectedCount        int     `json:"suspectedCount" bson:"suspectedCount"`
	CuredCount            int     `json:"curedCount" bson:"curedCount"`
	DeadCount             int     `json:"deadCount" bson:"deadCount"`
	Comment               string  `json:"comment" bson:"comment"`
	LocationID            int64   `json:"locationId" bson:"locationId"`
	Cities                []CityT `json:"cities" bson:"cities"`
	CountryName           string  `json:"countryName" bson:"countryName"`
	CountryEnglishName    string  `json:"countryEnglishName" bson:"countryEnglishName"`
	ContinentName         string  `json:"continentName" bson:"continentName"`
	ContinentEnglishName  string  `json:"continentEnglishName" bson:"continentEnglishName"`
	ProvinceEnglishName   string  `json:"provinceEnglishName" bson:"provinceEnglishName"`
	UpdateTime            int64   `json:"updateTime" bson:"updateTime"`
}

type CityT struct {
	CityName              string `json:"cityName" bson:"cityName"`
	CurrentConfirmedCount int    `json:"currentConfirmedCount" bson:"currentConfirmedCount"`
	ConfirmedCount        int    `json:"confirmedCount" bson:"confirmedCount"`
	SuspectedCount        int    `json:"suspectedCount" bson:"suspectedCount"`
	CuredCount            int    `json:"curedCount" bson:"curedCount"`
	DeadCount             int    `json:"deadCount" bson:"deadCount"`
	LocationID            int64  `json:"locationId" bson:"locationId"`
	CityEnglishName       string `json:"cityEnglishName" bson:"cityEnglishName"`
}

type AllProvinceDataT struct {
	Results []ProvinceT `json:"results" bson:"results"`
	Success bool        `json:"success" bson:"success"`
}

type AllProvinceNameT struct {
	Results []string `json:"results" bson:"results"`
	Success bool     `json:"success" bson:"success"`
}

func GetProvinceNames() (overall []string, err error) {
	url := fmt.Sprintf("https://lab.isaaclin.cn/nCoV/api/provinceName")

	body, err := HttpGet(url)

	result := AllProvinceNameT{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		ewlog.Error(err)
	}

	return result.Results, err
}

// dtypt: 1 最新，0 所有
func GetAllProvinceData(dtype int, province string) (all []ProvinceT, err error) {
	url := fmt.Sprintf("https://lab.isaaclin.cn/nCoV/api/area?latest=%d&province=%s", dtype, province)
	ewlog.Info(url)
	body, err := HttpGet(url)

	result := AllProvinceDataT{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		ewlog.Error(err)
	}

	return result.Results, err
}

// dtypt: 1 最新，0 所有
func RetryGetAllProvinceData(dtype int, province string, retryTimes int) (all []ProvinceT, err error) {
	for i := 0; i < retryTimes; i++ {
		all, err = GetAllProvinceData(dtype, province)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 20)
	}

	return
}
