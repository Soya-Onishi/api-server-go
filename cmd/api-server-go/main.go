package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Soya-Onishi/api-server-go/internal/controller"
	"github.com/Soya-Onishi/api-server-go/internal/repository"
	"github.com/gin-gonic/gin"
)

type DBProfile struct {
	user     string
	password string
	url      string
	dbname   string
}

var profile = DBProfile{
	user:     "app",
	password: "app",
	url:      "172.25.172.3:3306",
	dbname:   "todo",
}

func setupServer() *controller.Router {
	var db *sql.DB
	var err error

	dsn := fmt.Sprintf(
		"%v:%v@(%v)/%v",
		profile.user,
		profile.password,
		profile.url,
		profile.dbname,
	)
	db, err = sql.Open("mysql", dsn)
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
