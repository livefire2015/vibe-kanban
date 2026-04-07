package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers(mux *http.ServeMux, store *Store) {
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleList(w, r, store)
		case http.MethodPost:
			handleCreate(w, r, store)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPut:
			handleUpdate(w, r, store, id)
		case http.MethodDelete:
			handleDelete(w, store, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func handleList(w http.ResponseWriter, _ *http.Request, store *Store) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.List())
}

type createRequest struct {
	Title string `json:"title"`
}

var validStatuses = map[string]bool{
	"todo":        true,
	"in-progress": true,
	"done":        true,
}

func handleCreate(w http.ResponseWriter, r *http.Request, store *Store) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	task := store.Add(req.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

type updateRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

func handleUpdate(w http.ResponseWriter, r *http.Request, store *Store, id int) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if !validStatuses[req.Status] {
		http.Error(w, "status must be todo, in-progress, or done", http.StatusBadRequest)
		return
	}
	task, ok := store.Update(id, req.Title, req.Status)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func handleDelete(w http.ResponseWriter, store *Store, id int) {
	if !store.Delete(id) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
