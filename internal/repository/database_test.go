package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func createMockDB() (*Repository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	rep := NewRepository(db)

	return rep, mock
}

func TestGetAllTodos(t *testing.T) {
	rep, mock := createMockDB()

	todos := []TodoResponse{
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

	rows := mock.NewRows([]string{"id", "name"})
	for _, todo := range todos {
		rows.AddRow(todo.Id, todo.Name)
	}

	mock.ExpectQuery("SELECT \\* FROM todo_list ORDER BY id DESC").
		WillReturnRows(rows)

	actualTodos := rep.GetAllTodos()

	assert.NotNil(t, todos)
	assert.Equal(t, len(todos), len(actualTodos))
	for i := range todos {
		assert.Equal(t, todos[i].Id, actualTodos[i].Id)
		assert.Equal(t, todos[i].Name, actualTodos[i].Name)
	}
}
