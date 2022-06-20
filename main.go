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

	mux := newServeMux()

	mux.HandleFunc("/upload", fileSvc.upload)
	mux.HandleFunc("/show", fileSvc.show)
	mux.HandleFunc("/filenames", fileSvc.listFilenames)

	log.Println("listen on port: 7878")
	log.Println(http.ListenAndServe(":7878", mux))
}

type mux struct {
	*http.ServeMux
}

func newServeMux() *mux {
	return &mux{
		ServeMux: http.NewServeMux(),
	}
}

// TODO(quixote-liu): optimize the origin of username and password.
const (
	username = "admin@quixote_lcs"
	password = "hubei@lcs_1208"
)

func (mux *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// set cores headers
	setCorsHeaders(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// authenticate
	uname, pword, ok := r.BasicAuth()
	if !ok {
		responseJSON(w, http.StatusUnauthorized, H{
			"error": "authentication failed: missing username or password",
		})
		return
	}
	if uname != username || pword != password {
		responseJSON(w, http.StatusUnauthorized, H{
			"error": "authentication failed: username or password error",
		})
		return
	}

	mux.ServeMux.ServeHTTP(w, r)
}
