package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase(db *sql.DB) {

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
