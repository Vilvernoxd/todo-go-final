package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todo-go-final/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"}, http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusOK)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := dateOnly(now)

	if strings.TrimSpace(task.Date) == "" {
		task.Date = today.Format(DateLayout)
	}

	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return err
	}
	t = dateOnly(t)

	var next string
	if strings.TrimSpace(task.Repeat) != "" {
		next, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if t.Before(today) {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = today.Format(DateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data any, status ...int) {
	code := http.StatusOK
	if len(status) > 0 {
		code = status[0]
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
