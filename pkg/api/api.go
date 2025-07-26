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
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	log.Printf("Логируем параметры: now=%s, date=%s, repeat=%s", nowStr, dateStr, repeatStr)

	now, err := time.Parse(nextdate.DateFormatYMD, nowStr)
	if err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат даты"}, http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, dateStr, repeatStr)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обработка ошибки записи ответа
	if _, err := fmt.Fprintf(w, "%s", nextDate); err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
	}
}
