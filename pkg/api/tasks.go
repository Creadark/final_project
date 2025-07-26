package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/Creadark/final_project/pkg/db"
	"github.com/Creadark/final_project/pkg/nextdate"
)

const (
	limit = 10
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "Неправильный запрос"}, http.StatusMethodNotAllowed)
		return
	}

	// Получение параметров
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	var tasks []*db.Task
	var err error

	// Определение типа поиска
	if search != "" {
		parsedDate, err := time.Parse(nextdate.DateFormatDMYDot, search)
		if err == nil {
			// Преобразование в формат 20060102
			tasks, err = db.TasksByDate(parsedDate.Format(nextdate.DateFormatYMD), limit)
		} else {
			tasks, err = db.TasksBySearch(search, limit)
		}
	} else {
		// Обычный запрос без фильтра
		tasks, err = db.Tasks(limit)
	}

	// Обработка ошибок
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка получения задач"}, http.StatusInternalServerError)
		return
	}

	// Гарантируем пустой массив вместо null
	if tasks == nil {
		tasks = []*db.Task{}
	}

	// Формируем ответ с ключом "tasks"
	response := struct {
		Tasks []*db.Task `json:"tasks"`
	}{
		Tasks: tasks,
	}

	writeJson(w, response, http.StatusOK) // Один вызов writeJson
}
