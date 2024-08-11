package task

import (
	"net/http"

	"github.com/AramLab/todo-list/internal/service"
)

func RegisterRoutes(taskService service.TaskService, webDir string) http.Handler {
	mux := http.NewServeMux()

	h := NewHandlers(taskService)

	// Обработчик для `/api/nextdate`
	mux.HandleFunc("/api/nextdate", h.NextDateHandler)

	// Обработчик для `/api/tasks`
	mux.HandleFunc("/api/tasks", h.GetTasksHandler)

	// Обработчик для `/api/task` с разными методами
	mux.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTaskHandler(w, r)
		case http.MethodPost:
			h.AddTaskHandler(w, r)
		case http.MethodPut:
			h.PutTaskHandler(w, r)
		case http.MethodDelete:
			h.DeleteTaskHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для `/api/task/done`
	mux.HandleFunc("/api/task/done", h.DoneTaskHandler)

	// Обработчик для корневого URL, возвращающий `index.html`
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	return mux
}
