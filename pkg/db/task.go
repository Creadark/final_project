package db

import (
	"database/sql"
	"fmt"
)

// подключение к бд для инициализации в main.go
func InitDB(database *sql.DB) {
	db = database
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	// SQL-запрос с сортировкой по дате и лимитом
	query := `
        SELECT id, title, date, repeat 
        FROM scheduler 
        ORDER BY date ASC 
        LIMIT ?
    `

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к БД: %v", err)
	}
	defer rows.Close()

	var tasks []*Task

	// Итерация по результатам
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Date, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		tasks = append(tasks, &t)
	}

	// Проверка ошибок после итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки результатов: %v", err)
	}

	return tasks, nil
}

// TasksBySearch возвращает задачи по подстроке
func TasksBySearch(search string, limit int) ([]*Task, error) {
	query := `
        SELECT * FROM scheduler 
        WHERE title LIKE ? 
           OR comment LIKE ? 
        ORDER BY date 
        LIMIT ?
    `
	searchPattern := "%" + search + "%"
	return queryTasks(query, searchPattern, searchPattern, limit)
}

// TasksByDate возвращает задачи на конкретную дату
func TasksByDate(date string, limit int) ([]*Task, error) {
	query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
	return queryTasks(query, date, limit)
}

// queryTasks универсальная функция для выполнения запросов
func queryTasks(query string, args ...any) ([]*Task, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Date, &t.Repeat, &t.Comment); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %v", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки результатов: %v", err)
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE id = ?
    `
	var task Task
	err := db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	return &task, nil
}

// Обновление задачи
func UpdateTask(task *Task) error {
	query := `
        UPDATE scheduler 
        SET date = ?, title = ?, comment = ?, repeat = ? 
        WHERE id = ?
    `
	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("задача с id %s не существует", task.ID)
	}
	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %v", err)
	}
	return nil
}

// UpdateTaskDate обновляет дату выполнения задачи
func UpdateTaskDate(id, date string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err := db.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления: %v", err)
	}
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
