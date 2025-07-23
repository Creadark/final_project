package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormatYMD     = "20060102"   // 20060102
	DateFormatDMYDot  = "02.01.2006" // 02.01.2006
	DateFormatYMDDash = "2006-01-02" // 2006-01-02
)

// Проверка переданной даты относительно текущей
func afterNow(date, now time.Time) bool {
	date = date.Truncate(24 * time.Hour) // Убираем время, оставляя дату
	now = now.Truncate(24 * time.Hour)
	return date.After(now)
}

// Воскресенье становится 7, остальные дни сохраняют свои значения (1-6)
func weekdayToISO(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

// Возвращает следующую дату выполнения или ошибку
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Преобразует строку dstart в time.Time
	dateNow, err := time.Parse(DateFormatYMD, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты: %w", err)
	}
	// Если правило пустое - возвращает ошибку
	if len(repeat) == 0 {
		return "", fmt.Errorf("отсутствует правило")
	}
	// Разделяет строку repeat по пробелам
	repeatSplit := strings.Split(repeat, " ")
	rule := repeatSplit[0] // присваеваем значение первому переданному индексу в правиле

	// Проверяем первый индекс и выполняем нужные условия в зависимости от указанного в нем значения
	switch rule {
	case "y": // Добавляет 1 год, пока год не станет >= текущего
		for {
			// увеличиваем год на 1
			dateNow = dateNow.AddDate(1, 0, 0)

			// если дата равна или больше текущей, выходим
			if dateNow.Year() >= now.Year() {
				return dateNow.Format(DateFormatYMD), nil
			}
		}

	case "d": // Добавляет N дней, пока дата не станет позже текущей
		if len(repeatSplit) < 2 {
			return "", fmt.Errorf("нужен аргумент после d")
		}

		num, err := strconv.Atoi(repeatSplit[1])
		if err != nil {
			return "", fmt.Errorf("ошибка парсинга дней: %w", err)
		}

		if num < 1 || num > 400 {
			return "", fmt.Errorf("некорректный день")
		}

		for {

			dateNow = dateNow.AddDate(0, 0, num)
			if afterNow(dateNow, now) {
				break
			}
		}
		return dateNow.Format(DateFormatYMD), nil

	case "w": // Ищет ближайший день недели из списка, который будет после текущей даты
		if len(repeatSplit) < 2 {
			return "", fmt.Errorf("нужен аргумент после w")
		}

		days := []int{}
		numbers := strings.Split(repeatSplit[1], ",")
		for _, number := range numbers {
			day, err := strconv.Atoi(number)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("ошибка парсинга дня недели: %w", err)
			}
			days = append(days, day)
		}

		// ищем ближайший день недели
		for {
			if afterNow(dateNow, now) {
				currentDay := weekdayToISO(dateNow.Weekday())
				for _, d := range days {
					if d == currentDay {
						return dateNow.Format(DateFormatYMD), nil
					}
				}
			}
			dateNow = dateNow.AddDate(0, 0, 1)
		}

	case "m": // Ищет ближайший день месяца
		if len(repeatSplit) < 2 {
			return "", fmt.Errorf("не указан параметр")
		}

		// парсим дни
		dayParts := strings.Split(repeatSplit[1], ",")
		days := []int{}
		for _, el := range dayParts {
			day, err := strconv.Atoi(el)
			if err != nil || day == 0 || day < -31 || day > 31 {
				return "", fmt.Errorf("неверное значение дня")
			}
			days = append(days, day)
		}

		// парсим месяцы
		months := []int{}
		if len(repeatSplit) >= 3 {
			monthParts := strings.Split(repeatSplit[2], ",")
			for _, el := range monthParts {
				month, err := strconv.Atoi(el)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("неверное значение месяца")
				}
				months = append(months, month)
			}
		}

		for {
			if afterNow(dateNow, now) {
				year, month := dateNow.Year(), dateNow.Month()
				currentDay := dateNow.Day()

				// учитываем месяцы
				if len(months) > 0 {
					matchMonth := false
					for _, m := range months {
						if int(month) == m {
							matchMonth = true
							break
						}
					}
					if !matchMonth {
						dateNow = dateNow.AddDate(0, 0, 1)
						continue
					}
				}

				// учитываем дни
				lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
				for _, d := range days {
					if d == -3 {
						return "", fmt.Errorf("неверное значение")
					}
					targetDay := d
					if targetDay < 0 {
						targetDay = lastDayOfMonth + d + 1
					}
					if targetDay == currentDay {
						return dateNow.Format(DateFormatYMD), nil
					}
				}
			}
			dateNow = dateNow.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("некорректные параметры")
	}
}
