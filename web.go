package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var showDir bool
var sortByName bool
var files_pattern string

func (d webDir) Sort(a os.FileInfo, b os.FileInfo) (c bool) {
	if c = birthTime(a).After(birthTime(b)); sortByName {
		c = a.Name() < b.Name()
	}
	return
}

func (d webDir) Filter(fi os.FileInfo) bool {
	if fi.Name()[0] == 46 {
		return false
	} else if fi.IsDir() {
		return showDir
	} else if len(files_pattern) == 0 {
		return true
	}
	r, e := filepath.Match(files_pattern, fi.Name())
	if e != nil {
		return false
	}
	return r
}

func main() {

	portPtr := flag.Int("port", 8000, "port number")
	flag.BoolVar(&showDir, "dir", false, "show dir")
	flag.BoolVar(&sortByName, "sort", false, "sort by name")
	flag.StringVar(&files_pattern, "files", "", "files pattern")
	flag.Parse()

	cwd, _ := os.Getwd()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr), http.FileServer(webDir(cwd))))
}
