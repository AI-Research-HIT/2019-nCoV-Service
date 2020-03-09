package model

import (
	"fmt"
	"math"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/db"
	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/AI-Research-HIT/2019-nCoV-Service/util"
	"github.com/ender-wan/ewlog"
)

func Prediction(request protodef.PredictionRequest) (protodef.PredictionResponse, error) {
	var err error
	province := protodef.ProvinceData{}
	startDate, err := time.Parse(util.DateLayout, "2020-01-21")
	if err != nil {
		ewlog.Error(err)
		return protodef.PredictionResponse{}, err
	}
	val := PredictTemp{
		MList:     map[int]float64{0: 12.0, 5: 5.0, 10: 1.0},
		Infection: []int{},
	}

	if request.Template == 1 {
		key := fmt.Sprintf("%s-%d", request.Province, request.City)
		v, ok := TempPredictData[key]
		if ok {
			val.MList = v.MList
		}
	}

	if len(request.Mlist) > 0 {
		val.MList = request.Mlist
	}
	province, err = db.CalculateDataByDay(request.Province, request.City)
	if err != nil {
		ewlog.Error(err)
		return protodef.PredictionResponse{}, err
	}

	if len(province.Detail) < 6 {
		err = fmt.Errorf("%+v history to less: %d", request, len(province.Detail))
		ewlog.Error(err)
		return protodef.PredictionResponse{}, err
	}

	for _, d := range province.Detail {
		val.Infection = append(val.Infection, d.TotalInfection)
	}
	startDate, err = time.Parse(util.DateLayout, province.Detail[0].Date)
	if err != nil {
		ewlog.Error(err)
		return protodef.PredictionResponse{}, err
	}

	// 游离在外感染人群比例
	var alpha0 = 0.8
	// 感染者接触人数
	var m_list = val.MList
	// 潜伏人群与发病人群感染能力之比
	var k = 0.5
	cityData := []float64{}
	for _, data := range val.Infection {
		cityData = append(cityData, float64(data))
	}

	var Te = 4      //潜伏天数
	var beta = 0.05 //病毒本身感染力

	if request.Te != 0 {
		Te = request.Te
	}

	if request.Beta != 0.0 {
		beta = request.Beta
	}

	seir := SEIR(cityData, request.PredictDay, k, alpha0, Te, beta, m_list)
	I_list := []float64{}
	for i := 0; i < request.PredictDay+1; i++ {
		I_list = append(I_list, (1.0-alpha(float64(i), alpha0))*seir[1][i])
	}

	actives := []protodef.Active{}

	deltainfection := delta_Infection(seir)

	for i := 0; i < request.PredictDay; i++ {
		newInfection := 0
		totalInfection := 0
		if i < len(val.Infection) {
			totalInfection = val.Infection[i]
			if i == 0 {
				newInfection = val.Infection[i]
			} else {
				newInfection = val.Infection[i] - val.Infection[i-1]
			}
		}

		active := protodef.Active{
			NewInfection:   newInfection,
			TotalInfection: totalInfection,
			PredictNew:     int(deltainfection[i]),
			PredictTotal:   int(seir[2][i] + seir[3][i] + seir[4][i]),
			PredictRecover: int(seir[3][i]),
			PredictDeath:   int(seir[4][i]),
			Date:           startDate.AddDate(0, 0, i).Format("2006-01-02"),
			MVal:           seir[5][i],
			AlphaVal:       seir[6][i],
		}
		actives = append(actives, active)
	}

	resp := protodef.PredictionResponse{
		City:     request.City,
		Province: request.Province,
		Actives:  actives,
	}

	return resp, err
}

// alpha随时间递减
func alpha(x float64, alpha0 float64) float64 {
	return alpha0 - 0.01*x
}

func mOrg(x int, m_list []float64) float64 {
	if x < 5 {
		return m_list[0]
	} else if x >= 5 && x < 10 {
		return m_list[1]
	} else {
		return m_list[2]
	}
}

func r(x int) float64 {
	r_list := []float64{0.00000, 0.02, 0.04, 0.15}
	if x < 7 {
		return r_list[0]
	}
	if x >= 7 && x < 20 {
		return r_list[1]
	}

	return r_list[2]
}

func d(x int) float64 {
	d_list := []float64{0, 0.01, 0.015, 0.00, 0.00}
	if x < 7 {
		return d_list[0]
	}
	if x >= 7 && x < 12 {
		return d_list[1]
	}
	if x >= 12 && x < 13 {
		return d_list[2]
	}
	return d_list[3]
}

// 传播模型
// A 城市实际感染数据
// alpha0 游离在外感染人群比例
// k 潜伏人群与发病人群感染能力之比
// TT 传播天数
// m_list 感染者接触人数（分成三个档）
// Te 潜伏天数
// gamma 死亡率
// y 治愈率
// beta 病毒本身感染力
func SEIR(A []float64, TT int, k float64, alpha0 float64, Te int, beta float64, mlist map[int]float64) [][]float64 {
	var a = 1.0 / float64(Te)
	var E = []float64{}
	var I = []float64{}
	var R = []float64{}
	var Q = []float64{0.0}
	var D = []float64{}
	if Te <= 4 {
		E = append(E, A[5]-(4.0-float64(Te))*(A[5]-A[4])-A[0])
	} else {
		E = append(E, A[5]+(float64(Te)-4.0)*(A[6]-A[5])-A[0]) //潜伏者
	}
	I = append(I, (1.0/(1.0-alpha0))*float64(A[0])) //传染者
	R = append(R, 0)                                //恢复人数
	D = append(D, 0)
	Q[0] = A[0]

	alphaList := []float64{}
	mHistory := []float64{}
	m := 10.0
	for idx := 0; idx < TT; idx++ {
		fidx := float64(idx)
		alpha := alpha(fidx, alpha0)
		alphaList = append(alphaList, alpha)
		mval, ok := mlist[idx]
		if ok {
			m = mval
		}
		//m = float64(rand.Intn(10))
		mHistory = append(mHistory, m)

		e := E[idx] + alpha*m*I[idx]*beta - a*E[idx] + k*m*beta*E[idx]
		E = append(E, e)
		i := alpha*I[idx] + a*E[idx]
		I = append(I, i)
		q := Q[idx] + (1.0-alpha)*I[idx] - r(idx)*Q[idx] - d(idx)*Q[idx]
		Q = append(Q, q)
		r := R[idx] + r(idx)*Q[idx]
		R = append(R, r)
		d := D[idx] + d(idx)*Q[idx]
		D = append(D, d)
	}
	var SEIR_list = [][]float64{E, I, Q, R, D, mHistory, alphaList}
	return SEIR_list
}

func average(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, n := range nums {
		sum += n
	}

	return sum / float64(len(nums))
}

// 计算误差
func MAE(X []float64, X1 []float64) float64 {
	var ABS = []float64{}
	for i := 0; i < len(X); i++ {
		abs := math.Abs(X[i]-X1[i]) / X[i]
		ABS = append(ABS, abs)
	}

	return average(ABS)
}

// func SEIRPrediction(cityData []float64, TT int, k float64, m_list []float64, alpha0 float64) ([]float64, []float64) {
// 	I_list := []float64{}
// 	seir := SEIR(cityData, TT, k, m_list, alpha0)
// 	for i := 0; i < TT+1; i++ {
// 		I_list = append(I_list, (1.0-alpha(float64(i), alpha0))*seir[1][i])
// 	}
// 	rs := []float64{}
// 	for i, _ := range I_list {
// 		rs = append(rs, I_list[i]+seir[2][i]+seir[3][i])
// 	}

// 	deltainfection := delta_Infection(seir)

// 	return rs, deltainfection
// }

//每日新增感染数
func delta_Infection(seir [][]float64) []float64 {
	deltainfection := []float64{}
	predict_infection := []float64{}
	for i, _ := range seir[1] {
		predict_infection = append(predict_infection, seir[2][i]+seir[3][i]+seir[4][i])
	}

	for i := 1; i < len(seir[1]); i++ {
		deltainfection = append(deltainfection, predict_infection[i]-predict_infection[i-1])
	}

	return deltainfection
}
