package repository

import (
	"database/sql"
	"log"
	"os"
)

type Repository struct {
	db *sql.DB
}

type TodoResponse struct {
	Id   int
	Name string
}

func NewRepository(db *sql.DB) *Repository {
	r := new(Repository)
	r.db = db

	return r
}

func (r *Repository) GetAllTodos() []TodoResponse {
	rows, err := r.db.Query("SELECT * FROM todo_list ORDER BY id DESC")
	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return nil
	}
	defer rows.Close()

	resp := []TodoResponse{}
	for rows.Next() {
		var id int
		var name string

		rows.Scan(&id, &name)

		resp = append(resp, TodoResponse{
			Id:   id,
			Name: name,
		})
	}

	return resp
}
