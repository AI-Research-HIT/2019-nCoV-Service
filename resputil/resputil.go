package resputil

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ResponseData struct {
	Response

	Data interface{} `json:"data"`
}

func WriteFailed(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	resp := ResponseData{
		Response: Response{
			Code: code,
			Msg:  msg,
		},
	}
	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	} else {
		w.Write(body)
	}
}

func WriteSuccessWithData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := ResponseData{
		Response: Response{
			Code: 0,
			Msg:  "success",
		},
		Data: data,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	} else {
		w.Write(body)
	}
}

func WriteSuccessWithJsonData(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "application/json")
	resp := ResponseData{
		Response: Response{
			Code: 0,
			Msg:  "success",
		},
		Data: data,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	} else {
		w.Write(body)
	}
}

func WriteSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Code: 0,
		Msg:  "success",
	}

	body, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	} else {
		w.Write(body)
	}
}
