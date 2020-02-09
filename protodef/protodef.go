package protodef

type PredictionRequest struct {
	City string `json:"city"`
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
	City      string  `json:"city"`
	ActiveVal float64 `json:"activeVal`
	Date      int     `json:"date"`
}

type PredictionResponse struct {
	Virus   VirusResponse  `json:"virus`
	Actives ActiveResponse `json:"actives"`
}
