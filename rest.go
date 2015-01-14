package main

import (
	"net/http"
)

func APIReturn(code int, info string, w http.ResponseWriter) {
	output := make(map[string]map[string]interface{})
	output["result"] = make(map[string]interface{})
	output["result"]["code"] = code
	output["result"]["info"] = info
	WriteJSON(output, code, w)
}

func APISingleResult(resultCode int, resultInfo string, data map[string]interface{}, w http.ResponseWriter) {
	output := make(map[string]map[string]interface{})
	output["result"] = make(map[string]interface{})
	output["result"]["code"] = resultCode
	output["result"]["info"] = resultInfo
	output["data"] = data
	WriteJSON(output, resultCode, w)
}

type APIMultipleOutput struct {
	Result map[string]interface{}   `json:"result"`
	Data   []map[string]interface{} `json:"data"`
}

func APIMultipleResults(resultCode int, resultInfo string, data APIMultipleOutput, w http.ResponseWriter) {
	data.Result = make(map[string]interface{})
	data.Result["code"] = resultCode
	data.Result["info"] = resultInfo
	WriteJSON(data, resultCode, w)
}
