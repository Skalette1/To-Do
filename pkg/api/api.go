package api

import (
	"final-project/internal/auth"
	"final-project/pkg/utils"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", utils.NextDateHandler)
	http.HandleFunc("/api/task", auth.Middleware(TaskRouter))
	http.HandleFunc("/api/tasks", auth.Middleware(GetTasksHandler))
	http.HandleFunc("/api/task/done", auth.Middleware(DoneTaskHandler))
	http.HandleFunc("/api/signin", SignInHandler)
}

func TaskRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
