package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() {
	// Создание базы данных `scheduler.db`.
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Working directory", wd)
	db, err := sql.Open("sqlite3", wd+"/scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	defer func(d *sql.DB) {
		closeErr := d.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}(db)

	// Чтение файла `scheduler.sql`, в котором описана база данных.
	schema, err := os.ReadFile("scheduler.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Выполнение всех SQL-команд из файла `scheduler.sql`.
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
}
