package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer() http.Handler {
	store := newItemStore([]Item{
		{ID: "1", Name: "Item 1", Description: "Description 1"},
		{ID: "2", Name: "Item 2", Description: "Description 2"},
	})

	return newAPIServer(store).routes()
}

func TestGetItems(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/items", nil)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /items status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var items []Item
	if err := json.NewDecoder(recorder.Body).Decode(&items); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
}

func TestGetItemByID(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/items/1", nil)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /items/1 status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var item Item
	if err := json.NewDecoder(recorder.Body).Decode(&item); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if item.ID != "1" || item.Name != "Item 1" {
		t.Fatalf("item = %+v, want id 1 and name Item 1", item)
	}
}

func TestGetUnknownItemReturnsNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/items/unknown", nil)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("GET /items/unknown status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestCreateItem(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"Item 3","description":"Description 3"}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/items", body)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("POST /items status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	var item Item
	if err := json.NewDecoder(recorder.Body).Decode(&item); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if item.ID == "" || item.Name != "Item 3" {
		t.Fatalf("created item = %+v, want generated id and name Item 3", item)
	}
}

func TestUpdateItem(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"Updated","description":"Updated description"}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/items/1", body)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("PUT /items/1 status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var item Item
	if err := json.NewDecoder(recorder.Body).Decode(&item); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if item.ID != "1" || item.Name != "Updated" || item.Description != "Updated description" {
		t.Fatalf("updated item = %+v", item)
	}
}

func TestDeleteItem(t *testing.T) {
	server := testServer()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/items/1", nil)

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("DELETE /items/1 status = %d, want %d", recorder.Code, http.StatusNoContent)
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/items/1", nil)

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("GET /items/1 after delete status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestInvalidJSONReturnsBadRequest(t *testing.T) {
	body := bytes.NewBufferString(`{"name":`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/items", body)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("POST /items invalid JSON status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/items/1", nil)

	testServer().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PATCH /items/1 status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}
