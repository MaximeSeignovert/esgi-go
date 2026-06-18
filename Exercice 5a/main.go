package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type itemStore struct {
	mu    sync.RWMutex
	items []Item
}

func newItemStore(initialItems []Item) *itemStore {
	items := make([]Item, len(initialItems))
	copy(items, initialItems)

	return &itemStore{items: items}
}

func (s *itemStore) all() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Item, len(s.items))
	copy(items, s.items)
	return items
}

func (s *itemStore) findByID(id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.ID == id {
			return item, true
		}
	}

	return Item{}, false
}

func (s *itemStore) create(item Item) Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	item.ID = uuid.NewString()
	s.items = append(s.items, item)
	return item
}

func (s *itemStore) update(id string, update Item) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, item := range s.items {
		if item.ID == id {
			s.items[i].Name = update.Name
			s.items[i].Description = update.Description
			return s.items[i], true
		}
	}

	return Item{}, false
}

func (s *itemStore) delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, item := range s.items {
		if item.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}

	return false
}

type apiServer struct {
	store *itemStore
}

func newAPIServer(store *itemStore) *apiServer {
	return &apiServer{store: store}
}

func (a *apiServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/items", a.itemsHandler)
	mux.HandleFunc("/items/", a.itemByIDHandler)
	return mux
}

func (a *apiServer) itemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/items" {
		writeJSONError(w, http.StatusNotFound, "route not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.getItemsHandler(w, r)
	case http.MethodPost:
		a.createItemHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *apiServer) itemByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := itemIDFromPath(r.URL.Path)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.getItemHandler(w, r, id)
	case http.MethodPut:
		a.updateItemHandler(w, r, id)
	case http.MethodDelete:
		a.deleteItemHandler(w, r, id)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *apiServer) getItemsHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.store.all())
}

func (a *apiServer) getItemHandler(w http.ResponseWriter, _ *http.Request, id string) {
	item, found := a.store.findByID(id)
	if !found {
		writeJSONError(w, http.StatusNotFound, "item not found")
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (a *apiServer) createItemHandler(w http.ResponseWriter, r *http.Request) {
	item, ok := decodeItem(w, r)
	if !ok {
		return
	}

	created := a.store.create(item)
	writeJSON(w, http.StatusCreated, created)
}

func (a *apiServer) updateItemHandler(w http.ResponseWriter, r *http.Request, id string) {
	item, ok := decodeItem(w, r)
	if !ok {
		return
	}

	updated, found := a.store.update(id, item)
	if !found {
		writeJSONError(w, http.StatusNotFound, "item not found")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (a *apiServer) deleteItemHandler(w http.ResponseWriter, _ *http.Request, id string) {
	if !a.store.delete(id) {
		writeJSONError(w, http.StatusNotFound, "item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeItem(w http.ResponseWriter, r *http.Request) (Item, bool) {
	defer r.Body.Close()

	var item Item
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&item); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return Item{}, false
	}

	if strings.TrimSpace(item.Name) == "" {
		writeJSONError(w, http.StatusBadRequest, "name is required")
		return Item{}, false
	}

	return item, true
}

func itemIDFromPath(path string) (string, error) {
	id := strings.TrimPrefix(path, "/items/")
	if id == "" || id == path || strings.Contains(id, "/") {
		return "", errors.New("item id not found")
	}

	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encode JSON response: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func main() {
	store := newItemStore([]Item{
		{ID: "1", Name: "Clavier", Description: "Clavier mecanique compact"},
		{ID: "2", Name: "Souris", Description: "Souris sans fil ergonomique"},
	})
	server := newAPIServer(store)

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", server.routes()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
