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

type config struct {
	port         int
	showDir      bool
	sortByName   bool
	filesPattern string
}

var conf config = config{8000, false, false, ""}

func (d webDir) Sort(a os.FileInfo, b os.FileInfo) (c bool) {
	if c = birthTime(a).After(birthTime(b)); conf.sortByName {
		c = a.Name() < b.Name()
	}
	return
}

func (d webDir) Filter(fi os.FileInfo) bool {
	if fi.Name()[0] == 46 {
		return false
	} else if fi.IsDir() {
		return conf.showDir
	} else if len(conf.filesPattern) == 0 {
		return true
	}
	r, e := filepath.Match(conf.filesPattern, fi.Name())
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
		return []string{"-files", "*." + name}, nil
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
	return []string{"-files", "*." + name}, nil
}

func readArgs(args []string, c config) (config, string) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	port := f.Int("port", c.port, "port number")
	showDir := f.Bool("dir", c.showDir, "show dir")
	sortByName := f.Bool("sort", c.sortByName, "sort by name")
	filesPattern := f.String("files", c.filesPattern, "files pattern")
	f.Parse(args)
	params := ""
	if f.NArg() > 0 {
		params = f.Arg(0)
	}
	return config{*port, *showDir, *sortByName, *filesPattern}, params
}

func main() {

	var configName string
	conf, configName = readArgs(os.Args[1:], conf)
	if len(configName) > 0 {
		configLine, err := readConfig(configName)
		if err != nil {
			log.Fatal(err)
		}
		conf, _ = readArgs(configLine, conf)
		conf, _ = readArgs(os.Args[1:], conf)
	}

	cwd, _ := os.Getwd()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.port), http.FileServer(webDir(cwd))))
}
