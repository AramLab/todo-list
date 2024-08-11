package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DatePattern = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (t *Task) ValidateTitle() error {
	if t.Title == "" {
		return errors.New("не указан заголовок задачи")
	}
	return nil
}

func (t *Task) ValidateAndFormatDate() error {
	if t.Date == "" {
		t.Date = time.Now().Format(DatePattern)
	}

	date, err := time.Parse(DatePattern, t.Date)
	if err != nil {
		return errors.New("дата представлена в формате, отличном от 20060102")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if date.Before(today) {
		if t.Repeat == "" {
			t.Date = today.Format(DatePattern)
		} else {
			nextDate, err := t.NextDate(time.Now())
			if err != nil {
				return errors.New("правило повторения указано в неправильном формате")
			}
			t.Date = nextDate
		}
	}
	return nil
}

func (t *Task) ValidatePepeat() error {
	if t.Repeat != "" {
		valid := false
		switch t.Repeat[0] {
		case 'd':
			if len(t.Repeat) > 2 {
				days, err := strconv.Atoi(t.Repeat[2:]) // выводит подстроку начиная с индекса 2 и до конца.
				if err == nil && days > 0 && days <= 400 {
					valid = true
					return nil
				}
			}
		case 'y':
			if len(t.Repeat) == 1 {
				valid = true
				return nil
			}
		}
		if !valid {
			return errors.New("правило повторения указано в неправильном формате")
		}
	}
	return nil
}

// Реализует функцию повтора задачи.
func (t *Task) NextDate(now time.Time) (string, error) {
	taskDate, err := time.Parse(DatePattern, t.Date)
	if err != nil {
		return "", fmt.Errorf("неверная дата: %v", err)
	}

	if t.Repeat == "" {
		return "", errors.New("повторение не указано")
	}

	repeatParts := strings.Split(t.Repeat, " ")
	rule := repeatParts[0]

	switch rule {
	case "d":
		// Проверка на корректность формата вводимых данных.
		if len(repeatParts) < 2 {
			return "", errors.New("не указан интервал в днях")
		}

		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("недопустимое значение интервала в днях")
		}

		// Логика добавления дней.
		if taskDate.After(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}

		if taskDate.Before(now) {
			for taskDate.Before(now) {
				taskDate = taskDate.AddDate(0, 0, days)
			}
		}

	case "y":
		if taskDate.Year() < now.Year() {
			taskDate = taskDate.AddDate(now.Year()-taskDate.Year(), 0, 0)
		} else {
			taskDate = taskDate.AddDate(1, 0, 0)
		}

		if taskDate.Month() == time.February && taskDate.Day() == 29 && isLeapYear(taskDate.Year()) {
			taskDate = time.Date(taskDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
		}

	case "w", "m":
		return "", errors.New("неподдерживаемый формат")

	default:
		return "", errors.New("неподдерживаемый формат")
	}

	return taskDate.Format(DatePattern), nil
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
