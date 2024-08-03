package handlers

import (
	"net/http"
	"time"

	"github.com/AramLab/todo-list/utils"
)

/*
var webDir string = "./web"

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, webDir+"/index.html")
		return
	}
	http.ServeFile(w, r, webDir+r.URL.Path)
}
*/

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		now := r.URL.Query().Get("now")
		date := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")

		timeNow, err := time.Parse("20060102", now)
		if err != nil {
			return
		}

		nextDate, err := utils.NextDate(timeNow, date, repeat)
		if err != nil {
			return
		}

		response := nextDate
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}
