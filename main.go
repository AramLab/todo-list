package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	//"github.com/gorilla/mux"

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
	appPath, err := os.Executable()
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

	//r := mux.NewRouter()

	// Обработчик для `/api/nextdate`.
	http.HandleFunc("/api/nextdate", h.NextDateHandler) //.Methods("GET")

	// Обработчик для `/api/task`.
	//http.HandleFunc("/api/task", h.AddTaskHandler) //.Methods("POST")

	// Обработчик для `/api/tasks`.
	http.HandleFunc("/api/tasks", h.GetTasksHandler) //.Methods("GET")

	// Обработчик для `/api/task/{id}`.
	//http.HandleFunc("/api/task", h.GetTaskHandler) //.Methods("GET")

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTaskHandler(w, r)
		case http.MethodPost:
			h.AddTaskHandler(w, r)
		case http.MethodPut:
			h.PutTaskHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для корневого URL, возвращающий `index.html`.
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запуск сервера
	log.Printf("Starting server on :%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("ошибка запуска сервера")
		log.Fatal(err)
	}
}
