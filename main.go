package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AramLab/todo-list/database"
	"github.com/AramLab/todo-list/handlers"
)

var DB *sql.DB

func main() {

	// Открываем соединение с базой данных.
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	webDir := "./web"

	// Порт по умолчанию.
	port := "7540"

	// Записываем в `appPath` путь к исполняемому файлу.
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Записываем в `dbFile` путь, по которому должна храниться база данных.
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	// Проверяем наличие файла базы данных.
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		database.CreateDatabase(db)
	}

	// Создание экземпляра Handlers с подключением к базе данных
	h := handlers.NewHandlers(db)

	// Обработчик для `/api/nextdate`.
	http.HandleFunc("/api/nextdate", h.NextDateHandler)

	// Обработчик для `/api/tasks`.
	http.HandleFunc("/api/tasks", h.GetTasksHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Обработчик для `/api/task/{id}`.
			h.GetTaskHandler(w, r)
		case http.MethodPost:
			// Обработчик для `/api/task`.
			h.AddTaskHandler(w, r)
		case http.MethodPut:
			// Обработчик для `/api/task/{id}`.
			h.PutTaskHandler(w, r)
		case http.MethodDelete:
			// Обработчик для `/api/task/{id}`
			h.DeleteTaskHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для `/api/task/done`.
	http.HandleFunc("/api/task/done", h.DoneTaskHandler)

	// Обработчик для корневого URL, возвращающий `index.html`.
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запуск сервера
	log.Printf("Starting server on :%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
