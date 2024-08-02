package handlers

import (
	"net/http"
)

var webDir string = "./web"

func MainHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, webDir+"/index.html")
		return
	}
	http.ServeFile(w, r, webDir+r.URL.Path)
}
