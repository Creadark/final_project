package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Creadark/final_project/pkg/db"
	"github.com/Creadark/final_project/pkg/nextdate"
)

// writeJSON - общая функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id parameter is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": "failed to delete task"}, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "invalid repeat rule"}, http.StatusBadRequest)
			return
		}

		if err := db.UpdateTaskDate(id, nextDate); err != nil {
			writeJSON(w, map[string]string{"error": "failed to update task date"}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id parameter is required"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": "failed to delete task"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
