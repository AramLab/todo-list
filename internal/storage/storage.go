package storage

import "github.com/AramLab/todo-list/internal/entities"

type TaskRepository interface {
	Save(task *entities.Task) (int64, error)
	FindById(id string) (entities.Task, error)
	FindAll() ([]entities.Task, error)
	UpdateDate(id string, task *entities.Task) error
	Update(task *entities.Task) error
	DeleteById(id string) error
}

type Repository struct {
	TaskRepository TaskRepository
}
