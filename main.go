package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AramLab/todo-list/internal/server/http/task"
	domain "github.com/AramLab/todo-list/internal/service/domain/task"
	"github.com/AramLab/todo-list/internal/storage/database"
	"github.com/AramLab/todo-list/internal/storage/sqlite"
)

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

	repository := sqlite.NewRepository(db).TaskRepository
	service := domain.NewTaskService(repository)

	// Регистрация маршрутов
	router := task.RegisterRoutes(service, webDir)

	// Запуск сервера
	log.Printf("Starting server on :%s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
