package main

import (
	"net/http"
)

type webDir string

func (d webDir) Open(name string) (http.File, error) {
	return http.Dir(d).Open(name)
}
