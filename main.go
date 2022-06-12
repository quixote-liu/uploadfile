package main

import (
	"log"
	"net/http"
)

func main() {
	fileSvc := newFileService()
	if err := fileSvc.initDirs(); err != nil {
		log.Printf("initialize directories failed: %v", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", fileSvc.upload)

	log.Println("listen on port: 7878")
	log.Println(http.ListenAndServe(":7878", mux))
}
