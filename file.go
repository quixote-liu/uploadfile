package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const (
	dirNameVideo       string = "./video"
	dirNameImage       string = "./image"
	dirNameApplication string = "./application"
	dirNameText        string = "./text"
	dirNameAudio       string = "./audio"
	dirNameOther       string = "./other"

	dirContentTypeVideo       string = "video/*"
	dirContentTypeImage       string = "image/*"
	dirContentTypeApplication string = "application/*"
	dirContentTypeText        string = "text/*"
	dirContentTypeAudio       string = "audio/*"
	dirContentTypeOther       string = "*"
)

type fileService struct {
	subdirs  []*subDirectory
	otherDir *subDirectory
}

func newFileService() *fileService {
	return &fileService{
		subdirs: []*subDirectory{
			newSubDirectory(dirContentTypeVideo, dirNameVideo),
			newSubDirectory(dirContentTypeImage, dirNameImage),
			newSubDirectory(dirContentTypeApplication, dirNameApplication),
			newSubDirectory(dirContentTypeText, dirNameText),
			newSubDirectory(dirContentTypeAudio, dirNameAudio),
		},
		otherDir: newSubDirectory(dirContentTypeOther, dirNameOther),
	}
}

func (f *fileService) initDirs() error {
	for _, dir := range f.subdirs {
		if err := dir.make(); err != nil {
			return err
		}
	}
	if err := f.otherDir.make(); err != nil {
		return err
	}
	return nil
}

func (f *fileService) upload(w http.ResponseWriter, r *http.Request) {
	// check request method.
	if r.Method != http.MethodPost {
		responseJSON(w, http.StatusNotFound, H{"error": "the api only support POST method"})
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		responseJSON(w, http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	logs := newReadFileLogs()

loop:
	for {
		part, err := reader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break loop
			}
			logs.appendErr(err.Error())
			continue
		}

		var filename = part.FileName()
		subdirs := append(f.subdirs, f.otherDir)
		for _, dir := range subdirs {
			if dir.matched(filename) {
				if err := dir.writeFile(filename, part); err != nil {
					if errors.Is(err, errFileExist) {
						logs.appendFailed(filename, "Exist!")
						continue loop
					}
					log.Printf("failed: [%s]:%v", filename, err)
					logs.appendFailed(filename, err.Error())
					continue loop
				}
				logs.appendOK(filename)
				part.Close()
				continue loop
			}
		}
	}

	var status int
	if logs.onlyHasErrors() {
		status = http.StatusBadRequest
	} else {
		status = http.StatusCreated
	}
	responseJSON(w, status, logs.message())
}

func (f *fileService) show(w http.ResponseWriter, r *http.Request) {
	// check request method.
	if r.Method != http.MethodGet {
		responseJSON(w, http.StatusNotFound, H{"error": "the api only support GET method"})
		return
	}

	filename := r.URL.Query().Get("filename")

	var dir *subDirectory
	for _, d := range f.subdirs {
		if d.matched(filename) {
			dir = d
		}
	}
	if dir == nil {
		dir = f.otherDir
	}

	fs, err := dir.fs.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			responseJSON(w, http.StatusNotFound, H{"error": fmt.Sprintf("the file<%s> not found", filename)})
			return
		}
		log.Printf("[error]: open file<%s> failed: %v", filename, err)
		responseStatus(w, http.StatusNotFound)
		return
	}
	defer fs.Close()

	buf := bytes.NewBuffer(nil)
	defer buf.Reset()
	if _, err := buf.ReadFrom(fs); err != nil {
		log.Printf("read filename<%s> to buffer failed: %v", filename, err)
		responseStatus(w, http.StatusNotFound)
		return
	}

	ct := mime.TypeByExtension(filepath.Ext(filename))
	if ct == "" {
		ct = "text/plain, chartset=utf-8"
	}
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(200)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("response file content failed: %v", err)
		w.WriteHeader(500)
	}
}

func (f *fileService) listFilenames(w http.ResponseWriter, r *http.Request) {
	// check method: support GET.
	if r.Method != http.MethodGet {
		responseJSON(w, http.StatusNotFound, H{"error": "the api only support GET method"})
		return
	}

	filenames := []string{}
	errors := []string{}

	allSubDirs := append(f.subdirs, f.otherDir)
	for _, dir := range allSubDirs {
		entries, err := os.ReadDir(dir.path)
		if err != nil {
			errors = append(errors, err.Error())
		}
		for _, e := range entries {
			filenames = append(filenames, e.Name())
		}
	}

	var d H
	if len(errors) > 0 {
		d = H{"errors": errors, "filenames": filenames}
	} else {
		d = H{"filenames": filenames}
	}
	responseJSON(w, http.StatusOK, d)
}
