package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Repository struct {
	db *sql.DB
}

type TodoListManipulation interface {
	GetAllTodos() []TodoResponse
	PostTodo(todo TodoResponse) int
	DeleteTodo(id uint) int
	UpdateTodo(id int, todo TodoUpdater) int
}

type TodoResponse struct {
	Id   int
	Name string
}

type Updatable[T any] struct {
	Updatable bool
	Value     T
}
type TodoUpdater struct {
	Id   int
	Name Updatable[string]
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

func (r *Repository) UpdateTodo(id int, todo TodoUpdater) int {
	var namePart string
	if todo.Name.Updatable {
		namePart = fmt.Sprintf("title = '%v'", todo.Name.Value)
	} else {
		namePart = ""
	}

	doUpdate := false
	for _, update := range []bool{todo.Name.Updatable} {
		if update {
			doUpdate = true
		}
	}

	if doUpdate {
		sql := fmt.Sprintf("UPDATE todo.todo_list SET %v WHERE id = %v", namePart, id)
		return r.beginTx(func() error {
			if _, err := r.db.Exec(sql); err != nil {
				return err
			}

			return nil
		})
	} else {
		return http.StatusOK
	}

}
