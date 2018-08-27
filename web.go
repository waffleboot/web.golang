
package main

import (
	"log"
	"fmt"
	"flag"
	"strings"
	"path/filepath"
	"html/template"
	"os"
	"sort"
	"net/http"
	"strconv"
)

var showDir bool
var files_pattern string
var dir_html_template *template.Template

type dir_template_context struct {
	Name string
	Files []os.FileInfo
}

func (c dir_template_context) ShowFile(fi os.FileInfo) bool {
	if fi.Name()[0] == 46 {
		return false
	} else if fi.IsDir() {
		return showDir
	} else if len(files_pattern) == 0 {
		return true
	}
	r,e := filepath.Match(files_pattern,fi.Name())
	if e != nil { return false }
	return r
}

func show_dir(dir string, path string, resp http.ResponseWriter, req *http.Request) {
	file,e := os.Open(dir)
	if e != nil {
		http.NotFound(resp, req)
		return
	}
	defer file.Close()
	files,e := file.Readdir(0)
	if e != nil {
		http.Error(resp, fmt.Sprintf("unable to read content of dir %s", dir),http.StatusInternalServerError)
		return
	}
	sort.Slice(files,func(i,j int)bool{
		return files[i].ModTime().After(files[j].ModTime())
	})	
	if dir_html_template.Execute(resp,dir_template_context{path,files}) != nil {
		http.Error(resp, fmt.Sprintf("unable to show content of dir %s", dir),http.StatusInternalServerError)
		return
	}
}

func load_file(name string, resp http.ResponseWriter, req *http.Request) {
	log.Printf("load file %s",name)
	http.ServeFile(resp,req,name)
}

func make_dir_template() {
	var e error
	dir_html_template,e = template.New("dir").Parse(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN"><html>
	<title>Directory listing for {{.Name}}</title>
	<body>
	<h2>Directory listing for {{.Name}}</h2>
	<hr>
	<ul>
	{{$context := .}}
	{{range .Files}}
	{{if $context.ShowFile . }}<li><a href="{{.Name}}">{{.Name}}</a>{{end}}
	{{end}}
	</ul>
	<hr>
	</body>
	</html>`)
	if e != nil {
		log.Fatalf("unable to parse dir template")
	}
}

func main() {

	portPtr := flag.Int("port",8000,"port number")
	flag.BoolVar(&showDir,"dir",false,"show dir")
	flag.StringVar(&files_pattern,"files","","files pattern")
	flag.Parse()

	cwd,_ := os.Getwd()
	make_dir_template()
	http.HandleFunc("/",func (resp http.ResponseWriter, req *http.Request) {
		name := filepath.Join(cwd,filepath.FromSlash(req.URL.Path))
		fi,e := os.Stat(name)
		if e != nil {
			http.NotFound(resp, req)
			return
		}
		if fi.IsDir() {
			if strings.HasSuffix(req.URL.Path,"/") {
				show_dir(name,req.URL.Path,resp,req)
			} else {
				http.Redirect(resp,req,req.URL.Path+"/",http.StatusMovedPermanently)
			}
		} else {
			load_file(name,resp,req)
		}
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr),nil))
}