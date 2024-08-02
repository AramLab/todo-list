package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("неверная дата: %v", err)
	}

	if repeat == "" {
		return "", errors.New("повторение не указано")
	}

	repeatParts := strings.Split(repeat, "")
	rule := repeatParts[0]

	switch rule {
	case "d":
		if len(repeatParts) < 2 {
			return "", errors.New("не указан интервал в днях")
		}

		days, err := strconv.Atoi(repeatParts[1])

		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("недопустимое значение")
		}

		for taskDate.Before(now) || taskDate.Equal(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
	case "y":
		for taskDate.Before(now) || taskDate.Equal(now) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}

	case "w", "m":
		return "", errors.New("неподдерживаемый формат")
	}

	return taskDate.Format("20060102"), nil
}
