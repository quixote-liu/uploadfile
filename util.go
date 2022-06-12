package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type H map[string]interface{}

func responseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	v, err := json.Marshal(data)
	if err != nil {
		log.Printf("marshal data failed: %v", err)
		return
	}
	w.Write(v)
}
