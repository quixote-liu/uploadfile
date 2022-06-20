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

func responseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func setCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, "+r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
