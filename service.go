package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type ToDo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var DB *sql.DB

func InitDatabase() {
	var err error
	dbPath := "./awesome.db"
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to database at: %s", dbPath) // Log the database path
	_, err = DB.Exec(`
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        status TEXT
    );`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateToDo(title string, status string) (int64, error) {
	result, err := DB.Exec("INSERT INTO todos (title, status) VALUES (?, ?)", title, status)
	if err != nil {
		log.Printf("Failed to insert todo: %v", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		return 0, err
	}
	log.Printf("Inserted todo with ID: %d", id) // Log the inserted ID
	return id, nil
}

func DeleteToDo(id int64) error {
	result, err := DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no todo found with id %d", id)
	}
	return nil
}

func ReadToDoList() []ToDo {
	rows, err := DB.Query("SELECT id, title, status FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	todos := make([]ToDo, 0)
	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Status); err != nil {
			log.Fatal(err)
		}
		log.Printf("Fetched todo: %+v", todo) // Log each fetched todo
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return todos
}
