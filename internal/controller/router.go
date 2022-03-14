package controller

import (
	"net/http"
	"strconv"

	"github.com/Soya-Onishi/api-server-go/internal/repository"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	repo   *repository.Repository
}

func NewRouter(engine *gin.Engine, repo *repository.Repository) *Router {
	r := new(Router)
	r.engine = engine
	r.repo = repo

	r.setRouter(engine)

	return r
}

func (r *Router) Run() {
	r.engine.Run()
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Router) setRouter(e *gin.Engine) {
	e.GET("/", r.helloHandler)
	e.GET("/todos", r.returnTodo)
}

func (r *Router) helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func (r *Router) returnTodo(c *gin.Context) {
	todos := r.repo.GetAllTodos()
	if todos == nil {
		return
	}

	resp := []map[string]string{}
	for _, todo := range todos {
		m := map[string]string{
			"id":   strconv.Itoa(todo.Id),
			"name": todo.Name,
		}

		resp = append(resp, m)
	}

	c.JSON(http.StatusOK, resp)
}
