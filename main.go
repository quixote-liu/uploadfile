package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadFileHandler)

	log.Println(http.ListenAndServe(":7878", mux))
}

const MAX_UPLOAD_SIZE = 1024 * 1024 * 50 // 50MB

type H map[string]interface{}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// check request method.
	if r.Method != http.MethodPost {
		responseJSON(w, http.StatusBadRequest, H{"error": "the api only support post method"})
		return
	}

	// TODO(lcs): identity user.

	body := http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		responseJSON(w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

}

func responseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	v, err := json.Marshal(data)
	if err != nil {
		log.Println("marshal data failed: %v", err)
		return
	}
	w.Write(v)
}
