package model

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/AI-Research-HIT/2019-nCoV-Service/util"
	"github.com/ender-wan/ewlog"
)

type Person struct {
	PID             int64
	InfectorID      int64
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
	InfectPeople    []int64
	TreatmentDay    int
}

type PersonSlice []Person

func (s PersonSlice) Len() int { return len(s) }

func (s PersonSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s PersonSlice) Less(i, j int) bool { return s[i].PID < s[j].PID }

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
	BehaviorQuarantine
	BehaviorTreatment
)

func (p *Person) generateAge() {
	prob := rand.Float64()
	if prob < 0.17 {
		p.Age = rand.Intn(15)
	} else if prob < 0.9 {
		p.Age = rand.Intn(50) + 16
	} else {
		p.Age = rand.Intn(30) + 66
	}
}

func (p *Person) generateCureDay() {
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

func (p *Person) generateDeadProp() {
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

func (p *Person) generateTe() {
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

func (p *Person) generateM(baseM float64) {
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

func (p *Person) modifyTreatmentEffect(effect float64) {
	p.CuredDay = int(float64(p.CuredDay) * effect)
}

func (p *Person) generateBeta(baseBeta float64) {
	flucBeta := baseBeta / 10.0

	isPositive := rand.Intn(1)
	//persent := rand.Int31n(10) + 1
	//fluc := flucBeta / float64(persent)
	if isPositive == 0 {
		p.Beta = baseBeta + flucBeta
	} else {
		p.Beta = baseBeta - flucBeta
	}
}

func (p *Person) ChangeStatus() int {
	prob := rand.Float64()
	if prob > p.DeadProp {
		p.Status = StatusCured
	} else {
		p.Status = StatusDead
	}

	return p.Status
}

func (p *Person) InitInfectedPerson(baseM float64, baseBeta float64, infectDay int) {
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

func initInfectedSeed(num int, baseM float64, baseBeta float64, infectDay int, pid *int64) []Person {
	people := []Person{}
	for i := 0; i < num; i++ {
		p := Person{
			PID:          *pid,
			InfectorID:   0,
			InfectPeople: []int64{},
		}
		p.InitInfectedPerson(baseM, baseBeta, infectDay)
		people = append(people, p)
		*pid++
	}

	return people
}

func infectPerson(num int, baseM float64, baseBeta float64, infectDay int, treamentDay int, pid *int64, person *Person) []Person {
	people := []Person{}
	for i := 0; i < num; i++ {
		p := Person{
			PID:          *pid,
			InfectorID:   person.PID,
			InfectPeople: []int64{},
			TreatmentDay: treamentDay,
		}
		person.InfectPeople = append(person.InfectPeople, *pid)
		p.InitInfectedPerson(baseM, baseBeta, infectDay)
		people = append(people, p)
		*pid++
	}

	return people
}

func infectNum(p Person, idx int) int {
	if p.Beta < 0.0000001 {
		return 0
	}
	count := 0
	avg := p.Beta / float64(p.Te)
	weight := (idx - p.InfectStartDay) + 1
	realBeta := avg * float64(weight)
	for i := 0; i < int(p.M); i++ {
		prop := rand.Float64()
		if prop < realBeta {
			count++
		}
	}

	return count
}

func Simulate(seedNum, day int, mlist map[int]float64,
	betaList map[int]float64, TreatmentList map[int]int,
	treatmentEffect map[int]float64,
	startDate time.Time, isCalSpread,
	isQuarantine bool) (protodef.MonteCarloSimulationResp, error) {
	baseM, ok := mlist[0]
	if !ok {
		baseM = 2.0
	}

	var pid int64 = 1

	TreatmentTime, ok := TreatmentList[0]
	if !ok {
		TreatmentTime = 4
	}

	baseBeta, ok := betaList[0]
	if !ok {
		baseBeta = 0.05
	}

	teffect, ok := treatmentEffect[0]
	if !ok {
		teffect = 1.0
	}

	people := initInfectedSeed(seedNum, baseM, baseBeta, 0, &pid)

	statistics := []protodef.MonteCarloSimulationItem{protodef.MonteCarloSimulationItem{
		InfectedCount: seedNum,
		InfectedNew:   seedNum,
		Date:          startDate.Format(util.DateLayout)},
	}

	resp := protodef.MonteCarloSimulationResp{}

	startDate = startDate.AddDate(0, 0, 1)
	CuredPeople := []Person{}
	DeadPeople := []Person{}

	for i := 1; i < day; i++ {
		added := []Person{}

		stat := protodef.MonteCarloSimulationItem{}
		m, ok := mlist[i]
		if ok {
			baseM = m
		}

		tm, ok := TreatmentList[i]
		if ok {
			TreatmentTime = tm
		}

		beta, ok := betaList[i]
		if ok {
			baseBeta = beta
		}

		effect, ok := treatmentEffect[0]
		if ok {
			teffect = effect
		}

		for idx, p := range people {
			if p.Status == StatusCured || p.Status == StatusDead {
				continue
			}

			people[idx].generateM(baseM)
			people[idx].modifyTreatmentEffect(teffect)

			// 开始发病
			if p.InfectStartDay+p.Te == i {
				people[idx].Status = StatusAttacked
				people[idx].Behavior = BehaviorQuarantine
			}

			// 发现确诊并治疗
			if p.InfectStartDay+p.Te+p.TreatmentDay == i {
				people[idx].Behavior = BehaviorTreatment
				people[idx].ConfirmStartDay = i
				stat.ConfirmNew++

				if isQuarantine {
					for _, infector := range p.InfectPeople {
						for idx2, p2 := range people {
							if p2.PID == infector && p2.Behavior != BehaviorTreatment {
								people[idx2].Behavior = BehaviorQuarantine
							}
						}
					}
				}
			}

			if people[idx].Behavior == BehaviorFree &&
				(people[idx].Status == StatusInfected || people[idx].Status == StatusAttacked) {
				spreadNum := infectNum(p, i)
				if spreadNum > 0 {
					stat.InfectedNew += spreadNum

					new := infectPerson(spreadNum, baseM, baseBeta, i, TreatmentTime, &pid, &people[idx])
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

		tmpPeople := []Person{}

		for _, tp := range people {
			if tp.Status != StatusCured && tp.Status != StatusDead {
				tmpPeople = append(tmpPeople, tp)
			}
		}

		people = append(tmpPeople, added...)
		statistics = append(statistics, stat)
		//fmt.Println("累计: ", len(people))

		if len(people) > OneMillion {
			return resp, errors.New(fmt.Sprintf("在第%d天，%d人同时感染冠状病毒，局势已经失控, 停止仿真计算", i, len(people)))
		}
	}

	resp.Statistic = statistics

	if isCalSpread {
		totalPeople := map[int64]Person{}

		for _, p := range people {
			totalPeople[p.PID] = p
		}

		for _, p := range DeadPeople {
			totalPeople[p.PID] = p
		}

		for _, p := range CuredPeople {
			totalPeople[p.PID] = p
		}

		allPeople := []Person{}
		for _, v := range totalPeople {
			allPeople = append(allPeople, v)
		}

		sort.Sort(PersonSlice(allPeople))

		resp.SpreadTrack = statisticSpreadTrack(totalPeople, &allPeople)
	}
	//fmt.Println("总数", len(people))
	return resp, nil
}

func statisticSpreadTrack(totalPeople map[int64]Person, allPeople *[]Person) protodef.SpreadTrackResponse {
	spreadResp := protodef.SpreadTrackResponse{
		Nodes: []protodef.SpreadNode{},
		Links: []protodef.SpreadLink{},
	}

	count := 0

	addedPeople := map[int64]bool{}
	categories := []int64{}

	for _, v := range *allPeople {
		if count >= 2 {
			break
		}
		spreadLink(v, &addedPeople, &totalPeople, &spreadResp, v.PID, &categories)
		count++
	}
	spreadResp.Categories, spreadResp.Legends = generateCategories(&categories)

	return spreadResp
}

func generateCategories(catetories *[]int64) ([]protodef.SpreadCategory, []string) {
	cs := []protodef.SpreadCategory{}
	legends := []string{}

	for _, v := range *catetories {
		c := strconv.FormatInt(v, 10)
		cate := protodef.SpreadCategory{
			Name: c,
		}
		cs = append(cs, cate)
		legends = append(legends, c)
	}

	return cs, legends
}

func findCategoryIndex(pid int64, categories *[]int64) int {

	for i, v := range *categories {
		if v == pid {
			return i
		}
	}

	return -1
}

func spreadLink(person Person, added *map[int64]bool, totalPeople *map[int64]Person,
	resp *protodef.SpreadTrackResponse, pid int64, categories *[]int64) {
	if _, ok := (*added)[person.PID]; ok {
		return
	}

	if len(person.InfectPeople) > 0 {
		category := findCategoryIndex(pid, categories)
		if category < 0 {
			*categories = append(*categories, pid)
			category = len(*categories) - 1
		}

		(*resp).Nodes = append((*resp).Nodes, protodef.SpreadNode{
			Name:     strconv.FormatInt(person.PID, 10),
			Value:    len(person.InfectPeople),
			Category: category,
		})

		(*added)[person.PID] = true

		for _, val := range person.InfectPeople {
			(*resp).Links = append((*resp).Links, protodef.SpreadLink{
				Source: strconv.FormatInt(person.PID, 10),
				Target: strconv.FormatInt(val, 10),
			})
			if _, ok := (*added)[val]; ok {
				continue
			}
			p := (*totalPeople)[val]
			spreadLink(p, added, totalPeople, resp, pid, categories)
		}
	}
}

func printSpreadLevel(people map[int64]Person, pid int64, count int, spreadStat *[]int) {
	count++
	p, ok := people[pid]
	if ok {
		ewlog.Infof("pid %d infect %d people", pid, len(p.InfectPeople))
		if len(p.InfectPeople) > 0 {
			for _, ip := range p.InfectPeople {
				printSpreadLevel(people, ip, count, spreadStat)
			}
		}
	}
	*spreadStat = append(*spreadStat, count)
	ewlog.Infof("传播层数 %d", count)
}
