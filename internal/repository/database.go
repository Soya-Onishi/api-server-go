package repository

import (
	"database/sql"
	"log"
	"net/http"
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

func (r *Repository) beginTx(f func() error) int {
	tx, err := r.db.Begin()
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return http.StatusInternalServerError
	}

	if err := f(); err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (r *Repository) GetAllTodos() []TodoResponse {
	rows, err := r.db.Query("SELECT * FROM todo_list ORDER BY id")
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

func (r *Repository) PostTodo(todo TodoResponse) int {
	tx, err := r.db.Begin()
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return http.StatusInternalServerError
	}

	if _, err := tx.Exec("INSERT INTO todo.todo_list (title) VALUES(?)", todo.Name); err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (r *Repository) DeleteTodo(id uint) int {
	tx, err := r.db.Begin()
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return http.StatusInternalServerError
	}

	if _, err := r.db.Exec("DELETE FROM todo.todo_list WHERE id = ?", id); err != nil {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (r *Repository) UpdateTodo(id int, todo TodoResponse) int {
	return r.beginTx(func() error {
		if _, err := r.db.Exec("UPDATE todo.todo_list SET title = ? WHERE id = ?", todo.Name, id); err != nil {
			return err
		}

		return nil
	})
}
