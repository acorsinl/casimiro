package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func WriteJSON(input interface{}, status int, w http.ResponseWriter) {
	output, err := EncodeJSON(input)
	if err != nil {
		log.Println("JSON Encoding failed")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(output)
}

func DecodeJSON(r *http.Request, destination interface{}) error {
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, destination)
	if err != nil {
		return err
	}
	return nil
}

func EncodeJSON(input interface{}) ([]byte, error) {
	var b []byte
	var err error

	b, err = json.Marshal(input)
	if err != nil {
		log.Println("Encoding failed")
	}

	return b, nil
}
