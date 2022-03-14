package main

import (
	"database/sql"

	"github.com/Soya-Onishi/api-server-go/internal/controller"
	"github.com/Soya-Onishi/api-server-go/internal/repository"
	"github.com/gin-gonic/gin"
)

func setupServer() *controller.Router {
	var db *sql.DB
	var err error

	db, err = sql.Open("mysql", "")
	if err != nil {
		panic(err)
	}

	engine := gin.Default()
	repo := repository.NewRepository(db)

	return controller.NewRouter(engine, repo)
}

func main() {
	setupServer().Run()
}
