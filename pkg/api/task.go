package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Creadark/final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, struct {
			Error string `json:"error"`
		}{Error: "параметр id обязателен"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, struct {
			Error string `json:"error"`
		}{Error: err.Error()}, http.StatusNotFound)
		return
	}

	writeJson(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &task); err != nil {
		writeJson(w, struct {
			Error string `json:"error"`
		}{Error: "неверный JSON"}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJson(w, struct {
			Error string `json:"error"`
		}{Error: "id обязателен"}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJson(w, struct {
			Error string `json:"error"`
		}{Error: "title обязателен"}, http.StatusBadRequest)
		return
	}

	if err := verificationDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"status": "успешно обновлено"}, http.StatusOK)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r) // Перенесена проверка id внутрь getTaskHandler
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
