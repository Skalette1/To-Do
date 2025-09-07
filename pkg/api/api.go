package api

import (
	"final-project/internal/auth"
	"final-project/internal/metrics"
	"final-project/pkg/utils"
	"fmt"
	"net/http"
	"time"
)

func WithMetrics(name string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		h(rr, r)
		duration := time.Since(start).Seconds()
		metrics.IncRequests(name, r.Method, fmt.Sprint(rr.status))
		metrics.ObserveDuration(name, duration)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func Init() {
	http.HandleFunc("/api/nextdate", WithMetrics("nextdate", utils.NextDateHandler))
	http.HandleFunc("/api/task", auth.Middleware(WithMetrics("task", TaskRouter)))
	http.HandleFunc("/api/tasks", auth.Middleware(WithMetrics("tasks", GetTasksHandler)))
	http.HandleFunc("/api/task/done", auth.Middleware(WithMetrics("task_done", DoneTaskHandler)))
	http.HandleFunc("/api/signin", WithMetrics("signin", SignInHandler))
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
