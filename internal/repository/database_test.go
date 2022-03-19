package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var initDBData = []TodoResponse{
	{
		Id:   1,
		Name: "prepare hot water",
	},
	{
		Id:   2,
		Name: "wait for three minutes",
	},
	{
		Id:   3,
		Name: "eat ramen",
	},
}

func createRepository() *Repository {
	dsn := fmt.Sprintf("%v:%v@(%v)/%v", "app", "app", "db.test", "todo")
	txdb.Register("txdb", "mysql", dsn)

	db, err := sql.Open("txdb", uuid.New().String())
	if err != nil {
		panic(err)
	}

	repo := NewRepository(db)

	return repo
}

func TestGetAllTodos(t *testing.T) {
	rep := createRepository()
	defer rep.db.Close()

	expectTodos := initDBData
	actualTodos := rep.GetAllTodos()

	assert.Equal(t, len(expectTodos), len(actualTodos))
	for i := range expectTodos {
		assert.Equal(t, expectTodos[i].Id, actualTodos[i].Id)
		assert.Equal(t, expectTodos[i].Name, actualTodos[i].Name)
	}
}

func TestPostTodo(t *testing.T) {
	rep := createRepository()
	defer rep.db.Close()

	postTodos := []TodoResponse{
		{
			Id:   0,
			Name: "power on",
		},
		{
			Id:   0,
			Name: "erase directory",
		},
	}

	expectTodos := append(initDBData, postTodos...)

	for _, todo := range postTodos {
		status := rep.PostTodo(todo)
		assert.Equal(t, status, http.StatusOK)
	}

	actualTodos := rep.GetAllTodos()
	assert.NotNil(t, actualTodos)

	assert.Equal(t, len(expectTodos), len(actualTodos))
	for i := range expectTodos {
		assert.Equal(t, expectTodos[i].Name, actualTodos[i].Name)
	}
}

func TestDeleteTodo(t *testing.T) {
	t.Run("delete existance todo", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		status := rep.DeleteTodo(1)
		assert.Equal(t, status, http.StatusOK)

		expected := initDBData[1:]

		todos := rep.GetAllTodos()
		assert.NotNil(t, todos)
		assert.Equal(t, 2, len(todos))
		for i, todo := range todos {
			assert.Equal(t, expected[i].Name, todo.Name)
		}
	})
}
