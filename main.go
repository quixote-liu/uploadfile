package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileSvc := newFileService()
	if err := fileSvc.initDirs(); err != nil {
		log.Printf("initialize directories failed: %v", err)
		return
	}

	mux.HandleFunc("/upload", fileSvc.upload)

	log.Println(http.ListenAndServe(":7878", mux))
}
