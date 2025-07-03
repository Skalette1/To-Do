package api

import (
	"encoding/json"
	"final-project/pkg/db"
	"log"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method"))
		return
	}
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding body"))
		return
	}
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty title"))
		return
	}
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	if task.Repeat != "" {
		nDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error getting next date"))
			return
		}
		task.Date = nDate
	}
	taskTime, err := time.Parse("20060102", task.Date)
	if err != nil || taskTime.Before(now) {
		task.Date = now.Format("20060102")
	}
	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("error adding task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error adding task"))
		return
	}
	writeJson(w, id)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
