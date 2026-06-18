package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRouter() http.Handler {
	gin.SetMode(gin.TestMode)
	return NewAPI(NewTaskStore(initialTasks())).Router()
}

func TestGetTasks(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /tasks status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var tasks []Task
	if err := json.NewDecoder(recorder.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("len(tasks) = %d, want 2", len(tasks))
	}
}

func TestGetTask(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /tasks/1 status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var task Task
	if err := json.NewDecoder(recorder.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if task.ID != "1" || task.Title != "Apprendre Gin" {
		t.Fatalf("task = %+v, want id 1 and title Apprendre Gin", task)
	}
}

func TestGetUnknownTaskReturnsNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/tasks/unknown", nil)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("GET /tasks/unknown status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestCreateTaskRequiresAPIKey(t *testing.T) {
	body := bytes.NewBufferString(`{"title":"Nouvelle tache"}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", body)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("POST /tasks without key status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestCreateTask(t *testing.T) {
	body := bytes.NewBufferString(`{"title":"Nouvelle tache","description":"Depuis le test"}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", body)
	request.Header.Set("X-API-KEY", apiKey)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	var task Task
	if err := json.NewDecoder(recorder.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if task.ID == "" || task.Title != "Nouvelle tache" {
		t.Fatalf("created task = %+v, want generated id and title Nouvelle tache", task)
	}
}

func TestCreateTaskRequiresTitle(t *testing.T) {
	body := bytes.NewBufferString(`{"description":"Sans titre"}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", body)
	request.Header.Set("X-API-KEY", apiKey)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("POST /tasks without title status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUpdateTaskPartially(t *testing.T) {
	body := bytes.NewBufferString(`{"done":true}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/tasks/1", body)
	request.Header.Set("X-API-KEY", apiKey)

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("PUT /tasks/1 status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var task Task
	if err := json.NewDecoder(recorder.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if task.ID != "1" || task.Title != "Apprendre Gin" || !task.Done {
		t.Fatalf("updated task = %+v, want original title and done true", task)
	}
}

func TestDeleteTask(t *testing.T) {
	router := newTestRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	request.Header.Set("X-API-KEY", apiKey)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("DELETE /tasks/1 status = %d, want %d", recorder.Code, http.StatusNoContent)
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/tasks/1", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("GET /tasks/1 after delete status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}
