package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/Soya-Onishi/api-server-go/internal/controller"
	"github.com/Soya-Onishi/api-server-go/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	dsn := fmt.Sprintf("%v:%v@(%v)/%v", "app", "app", "db.test", "todo")
	txdb.Register("txdb", "mysql", dsn)
}

func setupMockServer() (*controller.Router, sqlmock.Sqlmock, error) {
	var db *sql.DB
	var err error

	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	engine := gin.Default()
	repo := repository.NewRepository(db)
	router := controller.NewRouter(engine, repo)

	return router, mock, nil
}

func TestHello(t *testing.T) {
	mockUserResp := `{"message":"Hello World"}`

	router, _, err := setupMockServer()
	assert.Nil(t, err)

	ts := httptest.NewServer(router.GetEngine())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, actual %v", resp.StatusCode)
	}

	respData, _ := ioutil.ReadAll(resp.Body)
	if string(respData) != mockUserResp {
		t.Fatalf("Expected response body %v, actual %v", mockUserResp, string(respData))
	}
}

func TestGetAllTodos(t *testing.T) {
	router, mock, err := setupMockServer()
	assert.Nil(t, err)

	mockTodoResp, _ := json.Marshal([]map[string]string{
		{
			"id":   "1",
			"name": "prepare hot water",
		},
		{
			"id":   "2",
			"name": "wait for three minutes",
		},
		{
			"id":   "3",
			"name": "eat ramen",
		},
	})

	rows := mock.NewRows([]string{"id", "name"}).
		AddRow(1, "prepare hot water").
		AddRow(2, "wait for three minutes").
		AddRow(3, "eat ramen")

	mock.ExpectQuery("SELECT .+ FROM todo_list ORDER BY id DESC").
		WillReturnRows(rows)

	ts := httptest.NewServer(router.GetEngine())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/todos", ts.URL))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, respBody, mockTodoResp)
}
