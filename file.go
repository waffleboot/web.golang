
package main

import ("net/http";"os")

type webDir string
type webFile struct {
	origin http.File
}

func (d webDir) Open(name string) (http.File, error) {
	f,e := http.Dir(d).Open(name)
	return webFile{f},e
}

func (f webFile) Read(p []byte) (int,error) {
	return f.origin.Read(p)
}

func (f webFile) Seek(offset int64, whence int) (int64, error) {
	return f.origin.Seek(offset, whence)
}

func (f webFile) Close() error {
	return f.origin.Close()
}

func (f webFile) Readdir(count int) ([]os.FileInfo, error) {
	fis,e := f.origin.Readdir(count)
	if e != nil {
		return fis,e
	}
	files := make([]os.FileInfo,0,len(fis))
	for _,fi := range fis {
		if showFile(fi) {
			files = append(files,fi)
		}
	}
	return files,e
}

func (f webFile) Stat() (os.FileInfo, error) {
	return f.origin.Stat()
}
