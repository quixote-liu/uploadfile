package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"regexp"
)

type subDirectory struct {
	contentType string
	path        string

	fs fs.FS
}

func newSubDirectory(contentType, path string) *subDirectory {
	return &subDirectory{
		contentType: contentType,
		path:        path,
	}
}

func (d *subDirectory) matched(filename string) bool {
	ext := mime.TypeByExtension(filepath.Ext(filename))
	match, _ := regexp.MatchString(d.contentType, ext)
	return match
}

func (d *subDirectory) make() error {
	if err := os.MkdirAll(d.path, fs.ModePerm); err != nil {
		return err
	}
	d.fs = os.DirFS(d.path)
	return nil
}

var errFileExist = errors.New("the file is exist")

func (d *subDirectory) writeFile(filename string, read io.Reader) error {
	fn := filepath.Join(d.path, filepath.Base(filename))
	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		return errFileExist
	}

	f, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("create file failed: %v", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, read); err != nil {
		return fmt.Errorf("write file content failed")
	}

	return nil
}
