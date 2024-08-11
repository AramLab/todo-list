package domain

import (
	"github.com/AramLab/todo-list/internal/service"
	"github.com/AramLab/todo-list/internal/service/domain/task"
	"github.com/AramLab/todo-list/internal/storage"
)

func NewService(repo *storage.Repository) *service.Service {
	return &service.Service{TaskService: task.NewTaskService(repo.TaskRepository)}
}
