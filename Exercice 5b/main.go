package main

import (
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const apiKey = "super-secret-key"

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TaskUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewTaskStore(initialTasks []Task) *TaskStore {
	tasks := make(map[string]Task, len(initialTasks))
	for _, task := range initialTasks {
		tasks[task.ID] = task
	}

	return &TaskStore{tasks: tasks}
}

func (s *TaskStore) All() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks
}

func (s *TaskStore) Find(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, found := s.tasks[id]
	return task, found
}

func (s *TaskStore) Create(task Task) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = uuid.NewString()
	s.tasks[task.ID] = task
	return task
}

func (s *TaskStore) Update(id string, update TaskUpdate) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, found := s.tasks[id]
	if !found {
		return Task{}, false
	}

	if update.Title != nil {
		task.Title = *update.Title
	}
	if update.Description != nil {
		task.Description = *update.Description
	}
	if update.Done != nil {
		task.Done = *update.Done
	}

	s.tasks[id] = task
	return task, true
}

func (s *TaskStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.tasks[id]; !found {
		return false
	}

	delete(s.tasks, id)
	return true
}

type API struct {
	store *TaskStore
}

func NewAPI(store *TaskStore) *API {
	return &API{store: store}
}

func (a *API) Router() *gin.Engine {
	router := gin.New()
	router.Use(LoggerMiddleware())

	router.GET("/tasks", a.GetTasks)
	router.GET("/tasks/:id", a.GetTask)

	protected := router.Group("/")
	protected.Use(AuthMiddleware(apiKey))
	protected.POST("/tasks", a.CreateTask)
	protected.PUT("/tasks/:id", a.UpdateTask)
	protected.DELETE("/tasks/:id", a.DeleteTask)

	return router
}

func (a *API) GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, a.store.All())
}

func (a *API) GetTask(c *gin.Context) {
	task, found := a.store.Find(c.Param("id"))
	if !found {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (a *API) CreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task JSON: title is required"})
		return
	}

	created := a.store.Create(task)
	c.JSON(http.StatusCreated, created)
}

func (a *API) UpdateTask(c *gin.Context) {
	var update TaskUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task JSON"})
		return
	}

	updated, found := a.store.Update(c.Param("id"), update)
	if !found {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "task not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (a *API) DeleteTask(c *gin.Context) {
	if !a.store.Delete(c.Param("id")) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "task not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Printf(
			"%s %s %s from %s in %s",
			start.Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			time.Since(start),
		)
	}
}

func AuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != expectedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or missing API key"})
			return
		}

		c.Next()
	}
}

func initialTasks() []Task {
	return []Task{
		{ID: "1", Title: "Apprendre Gin", Description: "Creer une API REST avec Gin", Done: false},
		{ID: "2", Title: "Tester les routes", Description: "Verifier les endpoints avec curl ou httptest", Done: true},
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	api := NewAPI(NewTaskStore(initialTasks()))
	log.Println("Server listening on http://localhost:8080")

	if err := api.Router().Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
