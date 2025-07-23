package api

import (
	"net/http"
	"time"

	"github.com/Creadark/final_project/pkg/db"
	"github.com/Creadark/final_project/pkg/nextdate"
)

// doneTaskHandler обрабатывает POST-запрос для завершения задачи
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "Неправильный запрос"}, http.StatusMethodNotAllowed)
		return
	}
	// Извлечение параметра id
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "параметр id обязателен"}, http.StatusBadRequest)
		return
	}

	// Получение задачи из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"}, http.StatusNotFound)
		return
	}

	// Обработка одноразовой задачи (repeat пуст)
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": "ошибка удаления задачи"}, http.StatusInternalServerError)
			return
		}
		writeJson(w, struct{}{}, http.StatusOK) // Пустой JSON {}
		return
	}

	// Для периодической задачи: расчёт следующей даты
	now := time.Now().UTC().Truncate(24 * time.Hour)
	nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "неверное правило повторения"}, http.StatusBadRequest)
		return
	}

	// Обновление даты задачи
	if err := db.UpdateTaskDate(id, nextDate); err != nil {
		writeJson(w, map[string]string{"error": "ошибка обновления даты"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJson(w, map[string]string{"error": "Неправильный запрос"}, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "параметр id обязателен"}, http.StatusBadRequest)
		return
	}

	// Проверка существования задачи
	_, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"}, http.StatusNotFound)
		return
	}

	// Удаление задачи
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": "ошибка удаления задачи"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK) // Пустой JSON {}
}
