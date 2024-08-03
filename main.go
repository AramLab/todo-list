package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AramLab/todo-list/database"
	"github.com/AramLab/todo-list/handlers"
)

func main() {

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
		database.CreateDatabase()
	}

	// Обработчик для `/api/nextdate`.
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)

	// Обработчик для корневого URL, возвращающий `index.html`.
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запуск сервера
	log.Printf("Starting server on :%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
