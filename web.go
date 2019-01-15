package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var showDir bool
var sortByName bool
var filesPattern string

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
	} else if len(filesPattern) == 0 {
		return true
	}
	r, e := filepath.Match(filesPattern, fi.Name())
	if e != nil {
		return false
	}
	return r
}

func readConfig(name string) ([]string, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.New("no current user")
	}
	p := filepath.Join(u.HomeDir, ".web.golang")
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.New("config file " + p + " not found")
	}
	r := bufio.NewReader(f)
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := strings.TrimSpace(s.Text())
		f := strings.Fields(t)
		if len(f) > 1 {
			k := f[0]
			if k == name {
				return f[1:], nil
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, errors.New("error on scanning " + p + " file")
	}
	return nil, errors.New("no entry for " + name)
}

func readArgs(args []string, d1 int, d2 bool, d3 bool, d4 string) (int, bool, bool, string, string) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	port := f.Int("port", d1, "port number")
	showDir := f.Bool("dir", d2, "show dir")
	sortByName := f.Bool("sort", d3, "sort by name")
	filesPattern := f.String("files", d4, "files pattern")
	f.Parse(args)
	params := ""
	if f.NArg() > 0 {
		params = f.Arg(0)
	}
	return *port, *showDir, *sortByName, *filesPattern, params
}

func main() {

	var port int
	var configName string
	port, showDir, sortByName, filesPattern, configName = readArgs(os.Args[1:], 8000, false, false, "")
	if len(configName) > 0 {
		configLine, err := readConfig(configName)
		if err != nil {
			log.Fatal(err)
		}
		port, showDir, sortByName, filesPattern, _ = readArgs(configLine, 8000, false, false, "")
		port, showDir, sortByName, filesPattern, _ = readArgs(os.Args[1:], port, showDir, sortByName, filesPattern)
	}

	cwd, _ := os.Getwd()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), http.FileServer(webDir(cwd))))
}
