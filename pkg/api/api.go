package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskRouter)
	http.HandleFunc("/api/tasks", TaskRouter) 
}

func TaskRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		TaskHandler(w, r)
	case http.MethodPost:
		AddTaskHandler(w, r)
	}
}
