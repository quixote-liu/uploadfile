package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

const MAX_UPLOAD_SIZE = 1024 * 1024 * 50 // 50MB

func (f *fileService) upload(w http.ResponseWriter, r *http.Request) {
	// check request method.
	if r.Method != http.MethodPost {
		responseJSON(w, http.StatusNotFound, H{"error": "the api only support post method"})
		return
	}

	readers, err := r.MultipartReader()
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, H{"error": err.Error()})
		return
	}

	var deliver = make(map[string]interface{})
	var index int
loop:
	for {
		index++
		indexStr := strconv.FormatInt(int64(index), 10)

		part, err := readers.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break loop
			}
			deliver[indexStr+"_error"] = err.Error()
			continue
		}

		var matched bool
		var contentType string = part.Header.Get("Content-Type")
		var filename = part.FileName()
		for _, dir := range f.subdirs {
			if dir.matched(filename, contentType) {
				if err := dir.writeFile(filename, part); err != nil {
					if errors.Is(err, errFileExist) {
						deliver[indexStr+"_error"] = fmt.Sprintf("the file<%s> is exist", filename)
						continue loop
					}
					deliver[indexStr+"_error"] = fmt.Sprintf("upload file<%s> failed: write file failed", filename)
					continue loop
				}
				matched = true
			}
		}
		if !matched {
			err := f.otherDir.writeFile(filename, part)
			if err != nil {
				deliver[indexStr+"_error"] = fmt.Sprintf("upload file<%s> failed: write file failed", filename)
				continue loop
			}
		}

		deliver[indexStr+"_success"] = fmt.Sprintf("<%s> uploaded", filename)
		part.Close()
	}

	responseJSON(w, http.StatusCreated, deliver)
}
