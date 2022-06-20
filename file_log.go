package main

import "strconv"

type readFileLogs struct {
	logs  map[string]string
	hasOK bool
}

func newReadFileLogs() *readFileLogs {
	return &readFileLogs{
		logs: make(map[string]string),
	}
}

func (d *readFileLogs) appendOK(filename string) {
	d.append(filename, "Uploaded Successful")
	d.hasOK = true
}

func (d *readFileLogs) appendErr(value string) {
	d.append("error", value)
}

func (d *readFileLogs) appendFailed(filename, value string) {
	d.append(filename, value)
}

func (d *readFileLogs) append(suffix, value string) {
	var index int
	if len(d.logs) == 0 {
		index = 1
	} else {
		index = len(d.logs) + 1
	}
	key := "[" + strconv.FormatInt(int64(index), 10) + "]" + suffix
	d.logs[key] = value
}

func (d *readFileLogs) onlyHasErrors() bool {
	return !d.hasOK
}

func (d *readFileLogs) message() map[string]string {
	return d.logs
}
