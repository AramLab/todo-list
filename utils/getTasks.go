package utils

import (
	"database/sql"
	//"log"

	"github.com/AramLab/todo-list/types"
)

func GetTasks(db *sql.DB) ([]types.Task, error) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		//log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []types.Task
	for rows.Next() {
		var task types.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			//log.Println("Error scanning row:", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		//log.Println("Error iterating over rows:", err)
		return nil, err
	}

	//log.Println("Retrieved tasks:", tasks)
	return tasks, nil
}
