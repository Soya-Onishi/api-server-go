package controller

import (
	"encoding/json"
	"errors"
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
	e.DELETE("/todos", r.deleteTodo)
	e.PATCH("/todos", r.updateTodo)
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

func errorHandling(err error, c *gin.Context) {
	log.SetOutput(os.Stderr)
	log.SetPrefix("[ERROR]")
	log.Printf("%v", err)

	c.AbortWithStatus(http.StatusBadRequest)
}

func (r *Router) postTodo(c *gin.Context) {
	var reqBody map[string]string
	reqBytes, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		errorHandling(err, c)
		return
	}

	if err := json.Unmarshal(reqBytes, &reqBody); err != nil {
		errorHandling(err, c)
		return
	}

	id, err := strconv.Atoi(reqBody["id"])
	if err != nil {
		errorHandling(err, c)
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

func getQueryID(c *gin.Context) (int, error) {
	idString, ok := c.GetQuery("id")
	if !ok {
		return -1, errors.New("query id does not exists")
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *Router) deleteTodo(c *gin.Context) {
	idString, ok := c.GetQuery("id")
	if !ok {
		errorHandling(errors.New("query id does not exists"), c)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		errorHandling(err, c)
		return
	}

	r.repo.DeleteTodo(uint(id))

	c.JSON(http.StatusOK, map[string]string{})
}

func (r *Router) updateTodo(c *gin.Context) {
	id, err := getQueryID(c)
	if err != nil {
		errorHandling(err, c)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	var body map[string]string
	json.Unmarshal(bodyBytes, &body)

	title, okTitle := body["name"]
	todo := repository.TodoUpdater{
		Id: id,
		Name: repository.Updatable[string]{
			Updatable: okTitle,
			Value:     title,
		},
	}

	r.repo.UpdateTodo(id, todo)
}
