package task

import (
	"database/sql"
	"errors"

	"github.com/AramLab/todo-list/internal/entities"
)

const limit = 50

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(task *entities.Task) (int64, error) {
	res, err := r.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
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

func (r *TaskRepository) FindById(id string) (entities.Task, error) {
	var task entities.Task
	row := r.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return entities.Task{}, errors.New(`{"error": "Задача не найдена"}`)
	}
	if err != nil {
		return entities.Task{}, errors.New(`{"error": "Ошибка при извлечении задачи"}`)
	}
	return task, nil
}

func (r *TaskRepository) FindAll() ([]entities.Task, error) {
	rows, err := r.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entities.Task
	for rows.Next() {
		var task entities.Task
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

func (r *TaskRepository) UpdateDate(id string, task *entities.Task) error {
	_, err := r.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", task.Date, id)
	if err != nil {
		return errors.New(`{"error":"Ошибка обновления даты задачи"}`)
	}
	return nil
}

func (r *TaskRepository) Update(task *entities.Task) error {
	_, err := r.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return errors.New(`{"error":"Ошибка обновления задачи"}`)
	}
	return nil
}

func (r *TaskRepository) DeleteById(id string) error {
	_, err := r.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return errors.New(`{"error": "Ошибка удаления задачи"}`)
	}
	return nil
}
