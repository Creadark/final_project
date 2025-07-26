package main

import (
	_ "modernc.org/sqlite"

	"log"
	"net/http"
	"os"

	"github.com/Creadark/final_project/pkg/api"
	"github.com/Creadark/final_project/pkg/db"
)

const (
	// DefaultPort - порт по умолчанию
	DefaultPort = "7540"
	// WebDir - директория с фронтендом
	WebDir = "./web"
)

// GetPort возвращает порт из переменной окружения или значение по умолчанию
func GetPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return DefaultPort
}

// SetupAndRun настраивает и запускает сервер
func SetupAndRun(port string) error {
	http.Handle("/", http.FileServer(http.Dir(WebDir)))

	log.Printf("Сервер запускается на порту %s", port) // Добавлено логирование
	return http.ListenAndServe(":"+port, nil)
}

func main() {
	api.Init()

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы %v", err)
	}
	defer db.Close()

	port := GetPort()
	if err := SetupAndRun(port); err != nil {
		log.Fatalf("Ошибка подключения сервера: %v", err)
	}

}
