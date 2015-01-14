package main

import (
	"code.google.com/p/go-uuid/uuid"
	"net/http"
	"time"
)

func Error(w http.ResponseWriter, error string, code int) {
	http.Error(w, error, code)
}

func UnixTimestamp() int32 {
	return int32(time.Now().Unix())
}

func NewUUID() string {
	return uuid.New()
}
