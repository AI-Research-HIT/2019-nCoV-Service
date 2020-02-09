package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AI-Research-HIT/2019-nCoV-Service/protodef"
	"github.com/ender-wan/ewlog"

	"github.com/AI-Research-HIT/2019-nCoV-Service/model"
	"github.com/AI-Research-HIT/2019-nCoV-Service/resputil"
)

func ModelCalculateHanlder(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ewlog.Error(err)
		resputil.WriteFailed(w, 2, err.Error())
		return
	}

	defer r.Body.Close()

	request := &protodef.PredictionRequest{}
	err = json.Unmarshal(data, request)

	virus, actives := model.Prediction(request.City)

	respData := protodef.PredictionResponse{
		Virus:   virus,
		Actives: actives,
	}
	resputil.WriteSuccessWithData(w, respData)
}
