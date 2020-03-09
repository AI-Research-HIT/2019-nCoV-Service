package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/AI-Research-HIT/2019-nCoV-Service/util"
)

type People struct {
	Status          int
	Behavior        int
	Te              int
	M               float64
	Beta            float64
	InfectStartDay  int
	DeadProp        float64
	Age             int
	ConfirmStartDay int
	CuredDay        int
}

const (
	StatusNormal = iota
	StatusInfected
	StatusAttacked
	StatusCured
	StatusDead
)

const OneMillion = 1000 * 1000

const (
	BehaviorFree = iota
	BehaviorHomeQuarantine
	BehaviorTreatment
)

func (p *People) generateAge() {
	prob := rand.Float64()
	if prob < 0.17 {
		p.Age = rand.Intn(15)
	} else if prob < 0.9 {
		p.Age = rand.Intn(50) + 16
	} else {
		p.Age = rand.Intn(30) + 66
	}
}

func (p *People) generateCureDay() {
	prob := rand.Float64()
	if prob < 0.1 {
		p.CuredDay = rand.Intn(3) + 5
	} else if prob < 0.3 {
		p.CuredDay = rand.Intn(8) + 2
	} else if prob < 0.7 {
		p.CuredDay = rand.Intn(10) + 5
	} else if prob < 0.9 {
		p.CuredDay = rand.Intn(15) + 4
	} else {
		p.CuredDay = rand.Intn(19) + 3
	}
}

func (p *People) generateDeadProp() {
	if p.Age > 65 {
		p.DeadProp = 0.05
	} else if p.Age >= 40 {
		p.DeadProp = 0.02
	} else if p.Age >= 14 {
		p.DeadProp = 0.01
	} else {
		p.DeadProp = 0.005
	}
}

func (p *People) generateTe() {
	prob := rand.Intn(100)
	if prob < 5 {
		p.Te = 14
	} else if prob < 10 {
		p.Te = 13
	} else if prob < 15 {
		p.Te = 12
	} else if prob < 30 {
		p.Te = rand.Intn(5) + 7
	} else if prob < 90 {
		p.Te = rand.Intn(4) + 3
	} else {
		p.Te = rand.Intn(2) + 1
	}
}

func (p *People) generateM(baseM float64) {
	if baseM < 3.0 {
		p.M = baseM
		return
	}
	flus := int(baseM) / 3
	isPositive := rand.Intn(1)
	//ewlog.Infof("flus: %d", flus)
	flus = rand.Intn(flus)
	if isPositive == 0 {
		p.M = baseM + float64(flus)
	} else {
		p.M = baseM - float64(flus)
	}

}

func (p *People) generateBeta(baseBeta float64) {
	flucBeta := baseBeta / 10.0

	isPositive := rand.Intn(1)
	persent := rand.Int31n(10) + 1
	fluc := flucBeta / float64(persent)
	if isPositive == 0 {
		p.Beta = baseBeta + fluc
	} else {
		p.Beta = baseBeta + fluc
	}
}

func (p *People) ChangeStatus() int {
	prob := rand.Float64()
	if prob > p.DeadProp {
		p.Status = StatusCured
	} else {
		p.Status = StatusDead
	}

	return p.Status
}

func (p *People) InitInfectedPerson(baseM float64, baseBeta float64, infectDay int) {
	p.Status = StatusInfected
	p.Behavior = BehaviorFree
	p.generateTe()
	p.generateM(baseM)
	p.generateBeta(baseBeta)
	p.InfectStartDay = infectDay
	p.generateAge()
	p.generateDeadProp()
	p.generateCureDay()
}

func initInfectedSeed(num int, baseM float64, baseBeta float64, infectDay int) []People {
	people := []People{}
	for i := 0; i < num; i++ {
		p := People{}
		p.InitInfectedPerson(baseM, baseBeta, infectDay)
		people = append(people, p)
	}

	return people
}

func infectPerson(num int, baseM float64, baseBeta float64, infectDay int) []People {
	people := []People{}
	for i := 0; i < num; i++ {
		p := People{}
		p.InitInfectedPerson(baseM, baseBeta, infectDay)
		people = append(people, p)
	}

	return people
}

func infectNum(p People, idx int) int {
	count := 0
	avg := p.Beta / float64(p.Te)
	weight := (idx - p.InfectStartDay) + 1
	realBeta := avg * float64(weight)
	for i := 0; i < int(p.M); i++ {
		prop := rand.Float64()
		if prop < realBeta {
			//ewlog.Infof("随机: %f, 感染率: %f, avg: %f, weight: %d, beta: %f, te: %d", prop, realBeta, avg, weight, p.Beta, p.Te)

			count++
		}
	}

	return count
}

func Simulate(seedNum, day int, mlist map[int]float64, baseBeta float64, TreatmentList map[int]int, startDate time.Time) ([]protodef.MonteCarloSimulationResponse, error) {
	baseM, ok := mlist[0]
	if !ok {
		baseM = 2.0
	}

	TreatmentTime, ok := TreatmentList[0]
	if ok {
		TreatmentTime = 4
	}

	people := initInfectedSeed(seedNum, baseM, baseBeta, 0)

	statistics := []protodef.MonteCarloSimulationResponse{protodef.MonteCarloSimulationResponse{
		InfectedCount: seedNum,
		InfectedNew:   seedNum,
		Date:          startDate.Format(util.DateLayout)},
	}
	startDate = startDate.AddDate(0, 0, 1)
	CuredPeople := []People{}
	DeadPeople := []People{}

	for i := 1; i < day; i++ {
		added := []People{}

		stat := protodef.MonteCarloSimulationResponse{}
		m, ok := mlist[i]
		if ok {
			baseM = m
		}

		tm, ok := TreatmentList[0]
		if ok {
			TreatmentTime = tm
		}

		for idx, p := range people {
			if p.Status == StatusCured || p.Status == StatusDead {
				continue
			}
			people[idx].generateM(baseM)
			if p.InfectStartDay+p.Te+TreatmentTime == i {
				people[idx].Behavior = BehaviorTreatment
				people[idx].ConfirmStartDay = i
				stat.ConfirmNew++
			}

			if p.InfectStartDay+p.Te == i {
				people[idx].Status = StatusAttacked
			}

			if people[idx].Behavior == BehaviorFree &&
				(people[idx].Status == StatusInfected || people[idx].Status == StatusAttacked) {
				spreadNum := infectNum(p, i)
				if spreadNum > 0 {
					stat.InfectedNew += spreadNum

					new := infectPerson(spreadNum, baseM, baseBeta, i)
					added = append(added, new...)
				}
			}

			if p.Behavior == BehaviorTreatment {
				if p.ConfirmStartDay+p.CuredDay == i {
					status := people[idx].ChangeStatus()
					if status == StatusCured {
						stat.CureNew++
						CuredPeople = append(CuredPeople, p)
					} else {
						stat.DeadNew++
						DeadPeople = append(DeadPeople, p)
					}

				}
			}
			if people[idx].Status == StatusInfected || people[idx].Status == StatusAttacked {
				stat.InfectingCount++
			}

			if people[idx].Behavior == BehaviorTreatment {
				stat.TreamentingCount++
			}
		}
		//fmt.Printf("Day %d 新增: %d \n", i, len(added))
		// fmt.Println("隔离人数: ", qCount)

		stat.InfectedCount = statistics[i-1].InfectedCount + stat.InfectedNew
		stat.ConfirmCount = statistics[i-1].ConfirmCount + stat.ConfirmNew
		stat.CureCount = statistics[i-1].CureCount + stat.CureNew
		stat.DeadCount = statistics[i-1].DeadCount + stat.DeadNew
		stat.Date = startDate.Format(util.DateLayout)
		startDate = startDate.AddDate(0, 0, 1)

		tmpPeople := []People{}

		for _, tp := range people {
			if tp.Status != StatusCured && tp.Status != StatusDead {
				tmpPeople = append(tmpPeople, tp)
			}
		}

		people = append(tmpPeople, added...)
		statistics = append(statistics, stat)
		//fmt.Println("累计: ", len(people))

		if len(people) > OneMillion {
			return statistics, errors.New(fmt.Sprintf("在第%d天，%d人同时感染冠状病毒，局势已经失控, 停止仿真计算", i, len(people)))
		}
	}

	//fmt.Println("总数", len(people))
	return statistics, nil
}
