package sqlite

import (
	"database/sql"

	"github.com/AramLab/todo-list/internal/storage"
	"github.com/AramLab/todo-list/internal/storage/sqlite/task"
)

func NewRepository(db *sql.DB) *storage.Repository {
	return &storage.Repository{TaskRepository: task.NewTaskRepository(db)}
}
