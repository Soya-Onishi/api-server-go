package repository

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"crypto/sha256"

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

func init() {
	dsn := fmt.Sprintf("%v:%v@(%v)/%v", "app", "app", "db.test", "todo")
	txdb.Register("txdb", "mysql", dsn)
}

func createRepository() *Repository {
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

	t.Run("delete non existance todo", func(t *testing.T) {
		repo := createRepository()
		defer repo.db.Close()

		status := repo.DeleteTodo(4)
		assert.Equal(t, http.StatusOK, status)

		todos := repo.GetAllTodos()
		assert.Equal(t, 3, len(todos))
		for i, todo := range todos {
			assert.Equal(t, initDBData[i].Name, todo.Name)
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	t.Run("update existance todo", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		updatedTitle := "title updated"
		status := rep.UpdateTodo(1, TodoUpdater{
			Id: 1,
			Name: Updatable[string]{
				Updatable: true,
				Value:     updatedTitle,
			},
		})

		assert.Equal(t, http.StatusOK, status)

		todos := rep.GetAllTodos()
		assert.NotNil(t, todos)
		assert.Equal(t, todos[0].Name, updatedTitle)

		assert.Equal(t, 3, len(todos))
		for i, todo := range todos[1:] {
			assert.Equal(t, initDBData[1:][i].Name, todo.Name)
		}
	})

	t.Run("update not existance todo", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		updateTitle := "title updated"
		status := rep.UpdateTodo(4, TodoUpdater{
			Id: 4,
			Name: Updatable[string]{
				Updatable: true,
				Value:     updateTitle,
			},
		})

		assert.Equal(t, http.StatusOK, status)

		todos := rep.GetAllTodos()
		for i, todo := range todos {
			assert.Equal(t, initDBData[i].Name, todo.Name)
		}
	})

	t.Run("no update cause no effect", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		status := rep.UpdateTodo(1, TodoUpdater{
			Id: 1,
			Name: Updatable[string]{
				Updatable: false,
				Value:     "",
			},
		})

		assert.Equal(t, http.StatusOK, status)

		todos := rep.GetAllTodos()
		for i, todo := range todos {
			assert.Equal(t, initDBData[i].Name, todo.Name)
		}
	})

	t.Run("update to empty string", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		status := rep.UpdateTodo(1, TodoUpdater{
			Id: 1,
			Name: Updatable[string]{
				Updatable: true,
				Value:     "",
			},
		})

		assert.Equal(t, http.StatusOK, status)

		todos := rep.GetAllTodos()
		assert.Equal(t, "", todos[0].Name)
	})
}

func TestGetUserInfo(t *testing.T) {
	t.Run("get user info by valid username", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		userInfo, status := rep.GetUserInfo("Taro")

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "Taro", userInfo.Username)
		assert.Equal(t, sha256.Sum256([]byte("Taro")), *userInfo.HashedPassword)
	})

	t.Run("get user info by invalid username", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		userInfo, status := rep.GetUserInfo("Unknown")

		assert.Nil(t, userInfo)
		assert.Equal(t, http.StatusUnauthorized, status)
	})
}

func TestGetSessionHash(t *testing.T) {
	t.Run("get session hash by valid username", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		hashSeed := fmt.Sprintf("%08x/%v", uint64(time.Now().Unix()), "Taro")
		hashExpect := sha256.Sum256([]byte(hashSeed))
		hash := fmt.Sprintf("%x", hashExpect)
		rep.db.Exec("UPDATE auth.users SET session_hash=? WHERE username='Taro'", hash)

		hashActual, status := rep.GetSessionHash("Taro")

		assert.NotNil(t, hashActual)
		assert.Equal(t, hashExpect, *hashActual)
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("get session but NULL", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		sessionHash, status := rep.GetSessionHash("Taro")

		assert.Nil(t, sessionHash)
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("get session hash by invalid username", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		sessionHash, status := rep.GetSessionHash("Unknown")

		assert.Nil(t, sessionHash)
		assert.Equal(t, http.StatusUnauthorized, status)
	})
}

func TestSetSessionHash(t *testing.T) {
	t.Run("set session hash by valid username", func(t *testing.T) {
		rep := createRepository()
		defer rep.db.Close()

		hash := sha256.Sum256([]byte{1, 2, 3})
		status := rep.SetSessionHash("Taro", hash)

		assert.Equal(t, http.StatusOK, status)

		rows, err := rep.db.Query("SELECT session_hash FROM auth.users WHERE username='Taro'")
		if err != nil {
			panic(err)
		}

		if !rows.Next() {
			panic("returns no rows from db")
		}

		var hashString string
		if err := rows.Scan(&hashString); err != nil {
			panic(err)
		}

		actualHash, _ := hex.DecodeString(hashString)
		assert.Equal(t, hash[:], actualHash)
	})
}
