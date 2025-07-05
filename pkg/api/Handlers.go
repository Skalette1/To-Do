package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"final-project/pkg/db"
	"final-project/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func getJWTSecret() []byte {
	if jwtSecret != nil {
		return jwtSecret
	}
	pass := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(pass))
	jwtSecret = hash[:]
	return jwtSecret
}

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
		nDate, err := utils.NextDate(now, task.Date, task.Repeat)
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
	toJson(w, id)
}

func toJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	tasks, err := db.GetAllTasks()
	if err != nil {
		log.Printf("error getting all tasks: %v", err)
		http.Error(w, "Error getting tasks", http.StatusInternalServerError)
		return
	}
	toJson(w, map[string]interface{}{"tasks": tasks})
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	toJson(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	var task db.Task
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &task); err != nil {
		http.Error(w, "Error decoding body", http.StatusBadRequest)
		return
	}

	if task.ID == 0 {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Empty title", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	if task.Repeat != "" {
		nDate, err := utils.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, "Error calculating next date", http.StatusInternalServerError)
			return
		}
		task.Date = nDate
	}
	taskTime, err := time.Parse("20060102", task.Date)
	if err != nil || taskTime.Before(now) {
		task.Date = now.Format("20060102")
	}

	if err := db.UpdateTask(&task); err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}
	toJson(w, map[string]string{"status": "updated"})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		http.Error(w, "Error deleting task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	toJson(w, map[string]string{"status": "deleted"})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	db.DoneTask(w, r)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || req.Password != pass {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный пароль"})
		return
	}

	now := time.Now().UTC()
	exp := now.Add(8 * time.Hour)

	hash := sha256.Sum256([]byte(pass))
	claims := jwt.MapClaims{
		"hash": hex.EncodeToString(hash[:]),
		"iat":  jwt.NewNumericDate(now),
		"exp":  jwt.NewNumericDate(exp),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(getJWTSecret())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка генерации токена"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
