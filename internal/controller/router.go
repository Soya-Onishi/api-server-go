package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Soya-Onishi/api-server-go/internal/repository"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	repo   repository.TodoListManipulation
}

func NewRouter(engine *gin.Engine, repo repository.TodoListManipulation) *Router {
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
	e.POST("/todos", r.postTodo)
}

func (r *Router) helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func (r *Router) returnTodo(c *gin.Context) {
	todos := r.repo.GetAllTodos()
	if todos == nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
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

func (r *Router) postTodo(c *gin.Context) {
	errorHandling := func(err error) {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[ERROR]")
		log.Printf("%v", err)

		c.AbortWithStatus(http.StatusBadRequest)
	}

	var reqBody map[string]string
	reqBytes, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		errorHandling(err)
		return
	}

	if err := json.Unmarshal(reqBytes, &reqBody); err != nil {
		errorHandling(err)
		return
	}

	id, err := strconv.Atoi(reqBody["id"])
	if err != nil {
		errorHandling(err)
		return
	}

	title, ok := reqBody["name"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	todo := repository.TodoResponse{
		Id:   id,
		Name: title,
	}

	status := r.repo.PostTodo(todo)

	c.JSON(status, map[string]string{})
}
