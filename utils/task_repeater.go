package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Реализует функцию повтора задачи.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("неверная дата: %v", err)
	}

	if repeat == "" {
		return "", errors.New("повторение не указано")
	}

	repeatParts := strings.Split(repeat, " ")
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

	return taskDate.Format("20060102"), nil
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
