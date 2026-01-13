package api

import (
	"net/http"

	"todo-go-final/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
