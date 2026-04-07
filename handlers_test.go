package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() http.Handler {
	store := NewStore()
	mux := http.NewServeMux()
	RegisterHandlers(mux, store)
	return mux
}

func TestCreateTask_Handler(t *testing.T) {
	srv := httptest.NewServer(setupRouter())
	defer srv.Close()

	body := bytes.NewBufferString(`{"title":"Test task"}`)
	resp, err := http.Post(srv.URL+"/api/tasks", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var task Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Title != "Test task" || task.Status != "todo" {
		t.Errorf("unexpected task: %+v", task)
	}
}

func TestListTasks_Handler(t *testing.T) {
	srv := httptest.NewServer(setupRouter())
	defer srv.Close()

	for _, title := range []string{"A", "B"} {
		body := bytes.NewBufferString(`{"title":"` + title + `"}`)
		http.Post(srv.URL+"/api/tasks", "application/json", body)
	}

	resp, err := http.Get(srv.URL + "/api/tasks")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var tasks []Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestUpdateTask_Handler(t *testing.T) {
	srv := httptest.NewServer(setupRouter())
	defer srv.Close()

	body := bytes.NewBufferString(`{"title":"Original"}`)
	http.Post(srv.URL+"/api/tasks", "application/json", body)

	body = bytes.NewBufferString(`{"title":"Updated","status":"done"}`)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/api/tasks/1", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var task Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Title != "Updated" || task.Status != "done" {
		t.Errorf("unexpected task: %+v", task)
	}
}

func TestDeleteTask_Handler(t *testing.T) {
	srv := httptest.NewServer(setupRouter())
	defer srv.Close()

	body := bytes.NewBufferString(`{"title":"To delete"}`)
	http.Post(srv.URL+"/api/tasks", "application/json", body)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/api/tasks/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}
