package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"final-project/pkg/db"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		data := map[string]interface{}{
			"error": err.Error(),
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonData)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(TaskResp{
		Tasks: tasks,
	})
}
