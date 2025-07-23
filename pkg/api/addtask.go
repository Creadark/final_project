package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Creadark/final_project/pkg/db"
	"github.com/Creadark/final_project/pkg/nextdate"
)

func verificationDate(task *db.Task) error {
	// Приводим текущее время в UTC и если не указана дата, ставим текущую
	now := time.Now().In(time.UTC)
	if task.Date == "" {
		task.Date = now.Format(nextdate.DateFormatYMD)
		log.Printf("Дата не указана, установлена текущая дата: %v", task.Date)
	}

	// Преобразовываем дату задачи

	t, err := time.Parse(nextdate.DateFormatYMD, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %v", err)
	}

	log.Printf("Парсинг даты задачи: %v", t)

	// Если указанная дата в прошлом, получаем следующую дату
	if !afterNow(now, t) {
		log.Printf("Ошибка, указанная дата (%v) меньше сегодняшней (%v)", t, now)

		// Если правило не указано, ставим текущую дату
		if task.Repeat == "" {
			task.Date = now.Format(nextdate.DateFormatYMD)
			log.Printf("Установлена текущая дата: %v", task.Date)
		} else {
			// Правило задано, указываем след. дату
			nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("некорректное правило повторения: %v", err)
			}

			// Преобразовывам к формату времени след. вычесленную дату
			parsedNextDate, err := time.Parse(nextdate.DateFormatYMD, nextDate)
			if err != nil {
				return fmt.Errorf("ошибка чтения следующей даты: %v", err)
			}

			log.Printf("Следующая дата: %v", parsedNextDate)

			// Повторение ежедневно
			if task.Repeat == "d 1" {
				// Устанавливаем дату задачи на сегодняшнюю, если она в будущем
				task.Date = now.Format(nextdate.DateFormatYMD)
				log.Printf("Ежедневное повторение, установлена сегодняшняя дата: %v", task.Date)
			} else if task.Repeat == "y" {
				// Повторение ежегодно
				// Если дата в прошлом, следующее повторение через год
				if !afterNow(now, parsedNextDate) {
					task.Date = parsedNextDate.AddDate(1, 0, 0).Format(nextdate.DateFormatYMD)
					log.Printf("Ежегодное повторение, дата перенесена на год: %v", task.Date)
				}
			} else {

				task.Date = nextDate
			}
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)}, http.StatusBadRequest)
		return
	}
	log.Printf("Полученные данные: %s", body)

	// Декодирование JSON
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Поле 'title' обязательно"}, http.StatusBadRequest)
		return
	}

	// Проверка даты
	if err := verificationDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка добавления задачи в базу данных: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Ответ с id задачи
	writeJson(w, map[string]int64{"id": id}, http.StatusOK)
}

func afterNow(now, date time.Time) bool {
	truncatedNow := now.Truncate(24 * time.Hour)                // Переводим в UTC и обрезаем время
	truncatedDate := date.In(time.UTC).Truncate(24 * time.Hour) // Переводим в UTC и обрезаем время

	log.Printf("Сравнение дат: now = %v, date = %v", truncatedNow, truncatedDate)
	return truncatedDate.After(truncatedNow) || truncatedDate.Equal(truncatedNow)
}
