package storage

import (
	"database/sql"
	"errors"

	"github.com/AramLab/todo-list/models"
)

// Get task.
func GetTask(id string, db *sql.DB) (models.Task, error) {
	var task models.Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return models.Task{}, errors.New(`{"error": "Задача не найдена"}`)
	}
	if err != nil {
		return models.Task{}, errors.New(`{"error": "Ошибка при извлечении задачи"}`)
	}
	return task, nil
}

const limit = 50

// Get tasks.
func GetTasks(db *sql.DB) ([]models.Task, error) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Delete.
func DeleteTaskById(id string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return errors.New(`{"error": "Ошибка удаления задачи"}`)
	}
	return nil
}

// Update.
func UpdateTask(task *models.Task, db *sql.DB) error {
	_, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return errors.New(`{"error":"Ошибка обновления задачи"}`)
	}
	return nil
}

// Update by id.
func UpdateTaskDate(id string, task *models.Task, db *sql.DB) error {
	_, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", task.Date, id)
	if err != nil {
		return errors.New(`{"error":"Ошибка обновления даты задачи"}`)
	}
	return nil
}

// Add.
func AddTask(task *models.Task, db *sql.DB) (int64, error) {
	if err := task.ValidateTitle(); err != nil {
		return 0, errors.New(`{"error":"title is required"}`)
	}

	// Проверка даты.
	if err := task.ValidateAndFormatDate(task); err != nil {
		return 0, errors.New(`{"error":"дата представлена в формате, отличном от 20060102"}`)
	}

	// Проверка шаблонов правил.
	if err := task.ValidatePepeat(); err != nil {
		return 0, errors.New(`{"error":"invalid repeat pattern"}`)
	}

	res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.New(`{"error":"Ошибка добавления задачи"}`)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New(`{"error":"Ошибка получения ID задачи"}`)
	}
	return id, nil
}
