package protodef

type PredictionRequest struct {
	City       int64           `json:"city"`
	Province   string          `json:"province"`
	Template   int             `json:"template"`
	PredictDay int             `json:"predictDay"`
	Te         int             `json:"te"`
	Beta       float64         `json:"beta"`
	Mlist      map[int]float64 `json:"mlist"`
}

type VirusResponse struct {
	City      string        `json:"city`
	VirusVals []VirusStatus `json:"virusVals"`
}

type VirusStatus struct {
	Cure      string `json:"cure"`
	Infection string `json:"infection"`
	Suspected string `json:"suspected"`
	Death     string `json:"death"`
}

type ActiveResponse struct {
	City       string   `json:"city`
	ActiveVals []Active `json:"activeVals`
}

type Active struct {
	NewInfection   int     `json:"newInfection"`
	TotalInfection int     `json:"totalInfection"`
	PredictNew     int     `json:"predictNew"`
	PredictTotal   int     `json:"predictTotal"`
	PredictRecover int     `json:"predictRecover"`
	PredictDeath   int     `json:"predictDeath"`
	MVal           float64 `json:"mval"`
	AlphaVal       float64 `json:"alpha"`
	Date           string  `json:"date"`
}

type PredictionResponse struct {
	City     int64    `json:"city`
	Province string   `json:"province"`
	Actives  []Active `json:"actives"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
}

type ProvinceData struct {
	PID          string       `json:"pid" bson:"pid"`
	ProvinceName string       `json:"provinceName" bson:"provinceName"`
	CityName     string       `json:"cityName" bson:"cityName"`
	Detail       []DetailData `json:"detail" bson:"detail"`
}

type DetailData struct {
	NewInfection   int    `json:"newInfection" bson:"newInfection"`
	NewDeath       int    `json:"newDeath" bson:"newDeath"`
	NewCure        int    `json:"newCure" bson:"newCure"`
	NewSusp        int    `json:"newSusp" bson:"newSusp"`
	Date           string `json:"date" bson:"date"`
	TotalInfection int    `json:"totalInfection" bson:"totalInfection"`
	TotalDeath     int    `json:"totalDeath" bson:"totalDeath"`
	TotalCure      int    `json:"totalCure" bson:"totalCure"`
	TotalSusp      int    `json:"totalSusp" bson:"totalSusp"`
}

type LatestData struct {
	NowInfection   int `json:"nowInfaction"`
	TotalCure      int `json:"totalCure"`
	TotalDeath     int `json:"totalDeath"`
	NowHeavy       int `json:"nowHeavy"`
	NowSusp        int `json:"nowSusp"`
	TotalInfection int `json:"totalInfection"`
}

type AllProviceDataRequest struct {
	ProvinceName string `json:"provinceName"`
	CityName     int64  `json:"cityName"`
}

type MonteCarloSimulationRequest struct {
	City               int64           `json:"city"`
	Province           string          `json:"province"`
	PredictDay         int             `json:"predictDay"`
	InitNum            int             `json:"initNum"`
	BetaList           map[int]float64 `json:"betalist"`
	Mlist              map[int]float64 `json:"mlist"`
	TreamentList       map[int]int     `json:"treamentList"`
	NoSymptomProb      float64         `json:"noSymptomProb"`
	SimulateNum        int             `json:"simNum"`
	IsQuarantineCloser bool            `json:"isQuarantineCloser"`
	MedicalNum         int             `json:"medicalNum"`
}

type MonteCarloSimulationResp struct {
	Statistic   []MonteCarloSimulationItem `json:"statistic"`
	SpreadTrack SpreadTrackResponse        `json:"spreadTrack"`
}

type MonteCarloSimulationItem struct {
	InfectedCount              int     `json:"infectedCount"`
	InfectedMin                int     `json:"infectedMin"`
	InfectedMax                int     `json:"infectedMax"`
	InfectedNew                int     `json:"infectedNew"`
	ConfirmCount               int     `json:"confirmCount"`
	ConfirmNew                 int     `json:"confirmNew"`
	RealConfirmCount           int     `json:"realConfirmCount"`
	RealConfirmNew             int     `json:"realConfirmNew"`
	DeadCount                  int     `json:"deadCount"`
	DeadNew                    int     `json:"deadNew"`
	RealDeadCount              int     `json:"realDeadCount"`
	RealDeadNew                int     `json:"realDeadNew"`
	CureCount                  int     `json:"cureCount"`
	CureNew                    int     `json:"cureNew"`
	RealCureCount              int     `json:"realCureCount"`
	RealCureNew                int     `json:"realCureNew"`
	Date                       string  `json:"date"`
	InfectingCount             int     `json:"infectingCount"`
	NoSymptomCount             int     `json:"noSymptomCount"`
	TreamentingCount           int     `json:"treamentingCount"`
	InfectedNotQuarantineCount int     `json:"infectedNotQuarantineCount"`
	CloserQuarantineCount      int     `json:"closerQuarantineCount"`
	AttackedNotTreamenting     int     `json:"attackedNotTreamenting"`
	ConfirmingCount            int     `json:"confirmingCount"`
	DeadProb                   float64 `json:"deadProb"`
}

type SpreadTrackResponse struct {
	Nodes      []SpreadNode     `json:"nodes"`
	Links      []SpreadLink     `json:"links"`
	Categories []SpreadCategory `json:"categories"`
	Legends    []string         `json:"legends"`
}

type SpreadNode struct {
	Name     string `json:"name"`
	Value    int    `json:"value"`
	Category int    `json:"category"`
}

type SpreadLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type SpreadCategory struct {
	Name string `json:"name"`
}

type CompareCountryResp struct {
	Countrys []string       `json:"countrys"`
	Data     []ProvinceData `json:"data"`
	Max      int            `json:"max"`
}
