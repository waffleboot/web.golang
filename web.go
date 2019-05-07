package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	g_port         = flag.Int("port", 8000, "port number")
	g_showDir      = flag.Bool("dir", false, "show dir")
	g_sortByName   = flag.Bool("sort", false, "sort by name")
	g_filesPattern = flag.String("files", "", "files pattern")
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
	} else if len(*g_filesPattern) == 0 {
		return true
	}
	r, e := filepath.Match(*g_filesPattern, fi.Name())
	if e != nil {
		return false
	}
	return r
}

func openConfigFile() *os.File {
	if currentUser, err := user.Current(); err == nil {
		file, _ := os.Open(filepath.Join(currentUser.HomeDir, ".web.golang"))
		return file
	}
	return nil
}

func readConfig(configKey string) []string {
	if configFile := openConfigFile(); configFile != nil {
		defer configFile.Close()
		r := bufio.NewReader(configFile)
		s := bufio.NewScanner(r)
		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			fields := strings.Fields(line)
			if len(fields) > 1 {
				key := fields[0]
				if key == configKey {
					return fields[1:]
				}
			}
		}
	}
	return []string{"-files", "*." + configKey}
}

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		flag.CommandLine.Parse(readConfig(flag.Arg(0)))
		flag.CommandLine.Parse(os.Args[1:])
	}
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
