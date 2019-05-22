package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var g_filesPattern []string

var (
	g_port       = flag.Int("port", 8000, "port number")
	g_showDir    = flag.Bool("dir", false, "show dir")
	g_sortByName = flag.Bool("sort", false, "sort by name")
)

func (d webDir) Sort(a os.FileInfo, b os.FileInfo) (c bool) {
	if c = birthTime(a).After(birthTime(b)); *g_sortByName {
		c = a.Name() < b.Name()
	}
	return
}

func (d webDir) Filter(fi os.FileInfo) bool {
	if fi.Name()[0] == 46 {
		return false
	} else if fi.IsDir() {
		return *g_showDir
	} else if len(g_filesPattern) == 0 {
		return true
	}
	for _, pattern := range g_filesPattern {
		if r, e := filepath.Match("*."+pattern, fi.Name()); e == nil && r {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	g_filesPattern = flag.Args()
	cwd, _ := os.Getwd()
	http.Handle("/", http.FileServer(webDir(cwd)))
	http.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		if cmd := exec.Command("pmset", "sleepnow"); cmd != nil {
			cmd.Run()
			os.Exit(0)
		}
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*g_port), nil))
}
