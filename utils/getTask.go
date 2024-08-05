package utils

import (
	"database/sql"
	"errors"

	"github.com/AramLab/todo-list/types"
)

func GetTask(id string, db *sql.DB) (types.Task, error) {
	var task types.Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return types.Task{}, errors.New(`{"error": "Задача не найдена"}`)
	} else if err != nil {
		return types.Task{}, errors.New(`{"error": "Ошибка при извлечении задачи"}`)
	} else {
		return task, nil
	}
}
