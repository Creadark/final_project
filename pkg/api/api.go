package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Creadark/final_project/pkg/nextdate"
)

// утилита для json ответов
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	// Логирование отправляемого ответа
	log.Printf("Response: %+v", data)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(w, "Ошибка сериализации ответа", http.StatusInternalServerError)
	}
}

// Обработчик для вычисления следующей даты
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	// Логируем параметры для отладки
	log.Printf("Логируем параметры: now=%s, date=%s, repeat=%s", nowStr, dateStr, repeatStr)

	// Преобразуем nowStr в формат time
	now, err := time.Parse(nextdate.DateFormatYMD, nowStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга now", http.StatusBadRequest)
		return
	}

	// Преобразуем dateStr в формат time
	date, err := time.Parse(nextdate.DateFormatYMD, dateStr)
	if err != nil {
		http.Error(w, "Ошибка парсинга date", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, date.Format(nextdate.DateFormatYMD), repeatStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка вычисления следующей даты: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", nextDate)
}
