package utils

import (
	"errors"
	"strconv"
)

func ValidateTitle(title string) error {
	if title == "" {
		return errors.New("не указан заголовок задачи")
	}
	return nil
}

func ValidatePepeat(repeat string) error {
	if repeat != "" {
		valid := false
		switch repeat[0] {
		case 'd':
			if len(repeat) > 2 {
				days, err := strconv.Atoi(repeat[2:]) // выводит подстроку начиная с индекса 2 и до конца.
				if err == nil && days > 0 && days <= 400 {
					valid = true
					return nil
				}
			}
		case 'y':
			if len(repeat) == 1 {
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
