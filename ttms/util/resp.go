package utils

import (
	"encoding/json"
	"net/http"
)

type H struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Rows    interface{} `json:"rows"`
	Total   interface{} `json:"total"`
}

func RespFail(writer http.ResponseWriter, err string) {
	Resp(writer, "-1", nil, err)
	return
}
func RespOk(writer http.ResponseWriter, date interface{}, msg string) {
	Resp(writer, "0", date, msg)
	return
}
func Resp(writer http.ResponseWriter, code string, data interface{}, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	h := H{
		Code:    code,
		Data:    data,
		Message: message,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	writer.Write(ret)
}
