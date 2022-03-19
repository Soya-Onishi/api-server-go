package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Soya-Onishi/api-server-go/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type RepositoryMock struct {
	todos []repository.TodoResponse
}

func (r *RepositoryMock) GetAllTodos() []repository.TodoResponse {
	return r.todos
}

func (r *RepositoryMock) PostTodo(todo repository.TodoResponse) int {
	r.todos = append(r.todos, todo)
	return http.StatusOK
}

func (r *RepositoryMock) DeleteTodo(id uint) int {
	var idx int = -1
	for i, todo := range r.todos {
		if uint(todo.Id) == id {
			idx = i
			break
		}
	}

	if idx != -1 {
		deleted := append(r.todos[:idx], r.todos[idx+1:]...)
		r.todos = deleted
	}

	return http.StatusOK
}

func (r *RepositoryMock) UpdateTodo(id int, todo repository.TodoResponse) int {
	var idx int = -1
	for i, todo := range r.todos {
		if todo.Id == id {
			idx = i
			break
		}
	}

	if idx != 1 {
		r.todos[idx] = todo
	}

	return http.StatusOK
}

var initDBData = []repository.TodoResponse{
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

func setupMock() *Router {
	mock := RepositoryMock{
		todos: initDBData,
	}

	return NewRouter(gin.Default(), &mock)
}

func TestRouteGetAllTodo(t *testing.T) {
	router := setupMock()
	ts := httptest.NewServer(router.engine)
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/todos", ts.URL))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respData []struct {
		Id   string
		Name string
	}

	respBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &respData)
	assert.Nil(t, err)
	assert.NotNil(t, respData)
	assert.Equal(t, 3, len(respData))
	for idx, todo := range respData {
		assert.Equal(t, initDBData[idx].Name, todo.Name)
	}
}
