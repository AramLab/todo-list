package service

import "github.com/AramLab/todo-list/internal/entities"

type TaskService interface {
	Create(task *entities.Task) (int64, error)
	Get(id string) (entities.Task, error)
	GetAll() ([]entities.Task, error)
	UpdateDate(id string, task *entities.Task) error
	Update(task *entities.Task) error
	Delete(id string) error
	Complete(id string) error
}

type Service struct {
	TaskService TaskService
}
