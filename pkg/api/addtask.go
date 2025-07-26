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
	now := time.Now().In(time.UTC)
	if task.Date == "" {
		task.Date = now.Format(nextdate.DateFormatYMD)
		log.Printf("Дата не указана, установлена текущая дата: %v", task.Date)
	}

	t, err := time.Parse(nextdate.DateFormatYMD, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %v", err)
	}

	log.Printf("Парсинг даты задачи: %v", t)

	if !afterNow(now, t) {
		log.Printf("Ошибка, указанная дата (%v) меньше сегодняшней (%v)", t, now)

		if task.Repeat == "" {
			task.Date = now.Format(nextdate.DateFormatYMD)
			log.Printf("Установлена текущая дата: %v", task.Date)
		} else {
			nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("некорректное правило повторения: %v", err)
			}

			parsedNextDate, err := time.Parse(nextdate.DateFormatYMD, nextDate)
			if err != nil {
				return fmt.Errorf("ошибка чтения следующей даты: %v", err)
			}

			log.Printf("Следующая дата: %v", parsedNextDate)

			switch task.Repeat {
			case "d 1":
				task.Date = now.Format(nextdate.DateFormatYMD)
				log.Printf("Ежедневное повторение, установлена сегодняшняя дата: %v", task.Date)
			case "y":
				if !afterNow(now, parsedNextDate) {
					task.Date = parsedNextDate.AddDate(1, 0, 0).Format(nextdate.DateFormatYMD)
					log.Printf("Ежегодное повторение, дата перенесена на год: %v", task.Date)
				}
			default:
				task.Date = nextDate
			}
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)}, http.StatusBadRequest)
		return
	}
	log.Printf("Полученные данные: %s", body)

	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Поле 'title' обязательно"}, http.StatusBadRequest)
		return
	}

	if err := verificationDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка добавления задачи в базу данных: %v", err)}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]int64{"id": id}, http.StatusOK)
}

func afterNow(now, date time.Time) bool {
	truncatedNow := now.Truncate(24 * time.Hour)
	truncatedDate := date.In(time.UTC).Truncate(24 * time.Hour)

	log.Printf("Сравнение дат: now = %v, date = %v", truncatedNow, truncatedDate)
	return truncatedDate.After(truncatedNow) || truncatedDate.Equal(truncatedNow)
}
