package main

import (
	"database/sql"
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
