package controller

import (
	"bytes"
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
	data := make([]repository.TodoResponse, len(initDBData))
	for i, todo := range initDBData {
		data[i].Id = todo.Id
		data[i].Name = todo.Name
	}

	mock := RepositoryMock{
		todos: data,
	}

	return NewRouter(gin.Default(), &mock)
}

func post(url string, content func() []byte) (*http.Response, error) {
	return nil, nil
}

func getTodo(ts *httptest.Server) []map[string]string {
	resp, err := http.Get(fmt.Sprintf("%v/todos", ts.URL))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var respData []map[string]string
	respBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &respData)
	if err != nil {
		panic(err)
	}

	return respData
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

func TestPostTodo(t *testing.T) {
	post := func(t *testing.T, ts *httptest.Server, body []byte) int {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%v/todos", ts.URL),
			bytes.NewBuffer(body),
		)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", http.DetectContentType(body))
		client := &http.Client{}
		resp, err := client.Do(req)

		assert.Nil(t, err)
		defer resp.Body.Close()

		return resp.StatusCode
	}

	todos := []map[string]string{
		{
			"id":   "4",
			"name": "new todo task",
		},
		{
			"id":   "4",
			"name": "another todo task",
		},
	}

	for _, todo := range todos {
		t.Run(todo["name"], func(t *testing.T) {
			router := setupMock()
			ts := httptest.NewServer(router.engine)
			defer ts.Close()

			body, err := json.Marshal(todo)
			if err != nil {
				panic(err)
			}

			assert.Equal(t, http.StatusOK, post(t, ts, body))

			getResp, err := http.Get(fmt.Sprintf("%v/todos", ts.URL))
			assert.Nil(t, err)
			defer getResp.Body.Close()

			var respData []struct {
				Id   string
				Name string
			}

			respBytes, _ := ioutil.ReadAll(getResp.Body)
			err = json.Unmarshal(respBytes, &respData)
			assert.Nil(t, err)
			assert.Equal(t, 4, len(respData))
			expect := append(initDBData, repository.TodoResponse{
				Id:   0,
				Name: todo["name"],
			})

			for idx, todo := range respData {
				assert.Equal(t, expect[idx].Name, todo.Name)
			}
		})
	}

	t.Run("without post todo", func(t *testing.T) {
		router := setupMock()
		ts := httptest.NewServer(router.engine)
		defer ts.Close()

		assert.Equal(t, http.StatusBadRequest, post(t, ts, make([]byte, 0)))
	})

	t.Run("invalid todo, todo title not existance", func(t *testing.T) {
		router := setupMock()
		ts := httptest.NewServer(router.engine)
		defer ts.Close()

		body, err := json.Marshal(map[string]string{
			"id": "4",
		})
		if err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusBadRequest, post(t, ts, body))
	})

	t.Run("invalid todo, todo id is invalid", func(t *testing.T) {
		router := setupMock()
		ts := httptest.NewServer(router.engine)
		defer ts.Close()

		body, err := json.Marshal(map[string]string{
			"id":   "abc",
			"name": "invalid todo",
		})
		if err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusBadRequest, post(t, ts, body))
	})
}

func TestDeleteTodo(t *testing.T) {
	runTest := func(f func(*httptest.Server)) {
		router := setupMock()
		ts := httptest.NewServer(router.engine)
		defer ts.Close()

		f(ts)
	}

	url := func(ts *httptest.Server) string {
		return fmt.Sprintf("%v/todos", ts.URL)
	}

	ids := []int{1, 2}

	for _, id := range ids {
		testTitle := fmt.Sprintf("delete todo via existance todo id[%v]", id)
		t.Run(testTitle, func(t *testing.T) {
			runTest(func(ts *httptest.Server) {
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("%v?id=%v", url(ts), id),
					bytes.NewBuffer(make([]byte, 0)),
				)
				if err != nil {
					panic(err)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				todos := getTodo(ts)
				assert.Equal(t, 2, len(todos))

				expect := make([]repository.TodoResponse, 0)
				expect = append(expect, initDBData[:id-1]...)
				expect = append(expect, initDBData[id:]...)
				for i, todo := range todos {
					assert.Equal(t, expect[i].Name, todo["name"])
				}
			})
		})
	}

	t.Run("Delete with no id cause error", func(t *testing.T) {
		runTest(func(ts *httptest.Server) {
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%v", url(ts)),
				bytes.NewBuffer(make([]byte, 0)),
			)
			if err != nil {
				panic(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("invalid id query, not number", func(t *testing.T) {
		runTest(func(ts *httptest.Server) {
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%v?id=abc", url(ts)),
				bytes.NewBuffer(make([]byte, 0)),
			)
			if err != nil {
				panic(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
