package db

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/cli"
	"github.com/AI-Research-HIT/2019-nCoV-Service/util"

	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
)

func CalculateDataByDay(provinceName string, city int64) (provinceData protodef.ProvinceData, err error) {
	result, err := FindAllProvinceData(provinceName)
	if err != nil {
		return
	}

	sort.Sort(util.ProvinceSlice(result))
	if city <= 0 {
		provinceData, err = CalculateProvinceDataByDay(result)
	} else {
		provinceData, err = CalculateCityDataByDay(result, city)
	}

	//ewlog.Infof("%+v", provinceData)

	return
}

func CalculateProvinceDataByDay(result []cli.ProvinceT) (provinceData protodef.ProvinceData, err error) {
	data := []cli.ProvinceT{}

	first := ""
	for i, d := range result {
		t := time.Unix(d.UpdateTime/1000, 0)

		dstr := t.Format(util.DateLayout)

		if strings.Compare(dstr, first) != 0 {
			data = append(data, d)
			first = dstr
		}

		if i == 0 {
			provinceData.ProvinceName = d.ProvinceName
			provinceData.PID = strconv.FormatInt(d.LocationID, 10)
			provinceData.Detail = []protodef.DetailData{}
		}
	}

	d0 := time.Unix(data[0].UpdateTime/1000, 0)
	d0Str := d0.Format(util.DateLayout)

	for i, d := range data {
		t := time.Unix(d.UpdateTime/1000, 0)

		dstr := t.Format(util.DateLayout)

		for strings.Compare(dstr, d0Str) != 0 {
			detail := protodef.DetailData{
				Date:           d0Str,
				TotalInfection: d.ConfirmedCount,
				TotalDeath:     d.DeadCount,
				TotalCure:      d.CuredCount,
				TotalSusp:      d.SuspectedCount,
			}
			if i > 0 {
				detail.NewInfection = data[i].ConfirmedCount - data[i-1].ConfirmedCount
				detail.NewDeath = data[i].DeadCount - data[i-1].DeadCount
				detail.NewCure = data[i].CuredCount - data[i-1].CuredCount
				detail.NewSusp = data[i].SuspectedCount - data[i-1].SuspectedCount
			}
			provinceData.Detail = append(provinceData.Detail, detail)
			d0 = d0.AddDate(0, 0, 1)
			d0Str = d0.Format(util.DateLayout)
		}

		detail := protodef.DetailData{
			Date:           d0Str,
			TotalInfection: d.ConfirmedCount,
			TotalDeath:     d.DeadCount,
			TotalCure:      d.CuredCount,
			TotalSusp:      d.SuspectedCount,
		}
		if i > 0 {
			detail.NewInfection = data[i].ConfirmedCount - data[i-1].ConfirmedCount
			detail.NewDeath = data[i].DeadCount - data[i-1].DeadCount
			detail.NewCure = data[i].CuredCount - data[i-1].CuredCount
			detail.NewSusp = data[i].SuspectedCount - data[i-1].SuspectedCount
		}
		provinceData.Detail = append(provinceData.Detail, detail)

		d0 = d0.AddDate(0, 0, 1)
		d0Str = d0.Format(util.DateLayout)
	}

	//ewlog.Infof("%+v", provinceData)

	return
}

func CalculateCityDataByDay(result []cli.ProvinceT, city int64) (provinceData protodef.ProvinceData, err error) {
	data := []cli.ProvinceT{}

	first := ""
	for i, d := range result {
		t := time.Unix(d.UpdateTime/1000, 0)

		dstr := t.Format(util.DateLayout)
		if strings.Compare(dstr, first) != 0 {
			data = append(data, d)
			first = dstr
		}

		if i == 0 {
			provinceData.ProvinceName = d.ProvinceName
			provinceData.PID = strconv.FormatInt(d.LocationID, 10)
			provinceData.Detail = []protodef.DetailData{}
		}
	}

	d0 := time.Unix(data[0].UpdateTime/1000, 0)
	d0Str := d0.Format(util.DateLayout)

	lastDetail := protodef.DetailData{}

	for _, d := range data {
		t := time.Unix(d.UpdateTime/1000, 0)

		dstr := t.Format(util.DateLayout)
		//ewlog.Info(dstr)

		for strings.Compare(d0Str, dstr) != 0 {
			c, ok := foundCityData(d.Cities, city)
			provinceData.CityName = c.CityName

			if ok {
				detail := protodef.DetailData{
					Date:           d0Str,
					TotalInfection: c.ConfirmedCount,
					TotalDeath:     c.DeadCount,
					TotalCure:      c.CuredCount,
					TotalSusp:      c.SuspectedCount,
				}
				// if i > 0 {
				// 	detail.NewInfection = data[i].ConfirmedCount - data[i-1].ConfirmedCount
				// 	detail.NewDeath = data[i].DeadCount - data[i-1].DeadCount
				// 	detail.NewCure = data[i].CuredCount - data[i-1].CuredCount
				// 	detail.NewSusp = data[i].SuspectedCount - data[i-1].SuspectedCount
				// }
				provinceData.Detail = append(provinceData.Detail, detail)
				lastDetail = detail
			} else {
				lastDetail.Date = d0Str
				provinceData.Detail = append(provinceData.Detail, lastDetail)
			}
			d0 = d0.AddDate(0, 0, 1)
			d0Str = d0.Format(util.DateLayout)
		}

		c, ok := foundCityData(d.Cities, city)
		provinceData.CityName = c.CityName

		if ok {
			detail := protodef.DetailData{
				Date:           d0Str,
				TotalInfection: c.ConfirmedCount,
				TotalDeath:     c.DeadCount,
				TotalCure:      c.CuredCount,
				TotalSusp:      c.SuspectedCount,
			}
			// if i > 0 {
			// 	detail.NewInfection = data[i].ConfirmedCount - data[i-1].ConfirmedCount
			// 	detail.NewDeath = data[i].DeadCount - data[i-1].DeadCount
			// 	detail.NewCure = data[i].CuredCount - data[i-1].CuredCount
			// 	detail.NewSusp = data[i].SuspectedCount - data[i-1].SuspectedCount
			// }
			provinceData.Detail = append(provinceData.Detail, detail)
			lastDetail = detail
		} else {
			lastDetail.Date = d0Str
			provinceData.Detail = append(provinceData.Detail, lastDetail)
		}
		d0 = d0.AddDate(0, 0, 1)
		d0Str = d0.Format(util.DateLayout)

	}

	// for _, d := range provinceData.Detail {

	// 	ewlog.Info(d.Date, " ", d.TotalInfection)
	// }

	//ewlog.Infof("%+v", provinceData)

	return
}

func foundCityData(citys []cli.CityT, locationId int64) (city cli.CityT, ok bool) {
	ok = false
	if len(citys) == 0 {
		return
	}

	for _, c := range citys {
		if c.LocationID == locationId {
			city = c
			ok = true
			return
		}
	}

	return
}
