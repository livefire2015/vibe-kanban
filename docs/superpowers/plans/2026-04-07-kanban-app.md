# Kanban App Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a kanban board with a Go REST API backend and SolidJS frontend.

**Architecture:** Go backend serves a REST API for CRUD operations on tasks stored in-memory. SolidJS frontend fetches tasks and renders them in three fixed columns (Todo, In Progress, Done). Backend also serves the built frontend static files.

**Tech Stack:** Go 1.21+ (standard library only), SolidJS, Vite, TypeScript

---

## File Structure

### Backend
- `main.go` — entry point, HTTP server setup, static file serving
- `handlers.go` — HTTP handler functions for task CRUD
- `store.go` — in-memory task store with mutex
- `handlers_test.go` — tests for all endpoints

### Frontend
- `frontend/package.json` — dependencies
- `frontend/tsconfig.json` — TypeScript config
- `frontend/vite.config.ts` — Vite config with proxy for dev
- `frontend/index.html` — HTML entry point
- `frontend/src/index.tsx` — app entry, render root
- `frontend/src/App.tsx` — main app component
- `frontend/src/App.css` — styles

---

### Task 1: Go Backend — Store

**Files:**
- Create: `store.go`
- Create: `store_test.go`

- [ ] **Step 1: Write failing tests for the store**

```go
// store_test.go
package main

import "testing"

func TestAddTask(t *testing.T) {
	s := NewStore()
	task := s.Add("Buy groceries")
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.Title != "Buy groceries" {
		t.Errorf("expected title 'Buy groceries', got %s", task.Title)
	}
	if task.Status != "todo" {
		t.Errorf("expected status 'todo', got %s", task.Status)
	}
}

func TestListTasks(t *testing.T) {
	s := NewStore()
	s.Add("Task 1")
	s.Add("Task 2")
	tasks := s.List()
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	s := NewStore()
	s.Add("Old title")
	updated, ok := s.Update(1, "New title", "in-progress")
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if updated.Title != "New title" || updated.Status != "in-progress" {
		t.Errorf("unexpected task: %+v", updated)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	s := NewStore()
	_, ok := s.Update(99, "x", "todo")
	if ok {
		t.Fatal("expected update to fail for missing task")
	}
}

func TestDeleteTask(t *testing.T) {
	s := NewStore()
	s.Add("To delete")
	ok := s.Delete(1)
	if !ok {
		t.Fatal("expected delete to succeed")
	}
	if len(s.List()) != 0 {
		t.Fatal("expected 0 tasks after delete")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	s := NewStore()
	ok := s.Delete(99)
	if ok {
		t.Fatal("expected delete to fail for missing task")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -v -run TestAdd`
Expected: FAIL — `NewStore` not defined

- [ ] **Step 3: Implement the store**

```go
// store.go
package main

import "sync"

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type Store struct {
	mu     sync.Mutex
	tasks  []Task
	nextID int
}

func NewStore() *Store {
	return &Store{nextID: 1}
}

func (s *Store) Add(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Task{ID: s.nextID, Title: title, Status: "todo"}
	s.nextID++
	s.tasks = append(s.tasks, t)
	return t
}

func (s *Store) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

func (s *Store) Update(id int, title, status string) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Title = title
			s.tasks[i].Status = status
			return s.tasks[i], true
		}
	}
	return Task{}, false
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -v`
Expected: All 6 tests PASS

- [ ] **Step 5: Commit**

```bash
git add store.go store_test.go
git commit -m "feat: add in-memory task store with tests"
```

---

### Task 2: Go Backend — HTTP Handlers

**Files:**
- Create: `handlers.go`
- Create: `handlers_test.go`

- [ ] **Step 1: Write failing tests for handlers**

```go
// handlers_test.go
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

func TestCreateTask(t *testing.T) {
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

	// Create two tasks
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

	// Create
	body := bytes.NewBufferString(`{"title":"Original"}`)
	http.Post(srv.URL+"/api/tasks", "application/json", body)

	// Update
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

	// Create
	body := bytes.NewBufferString(`{"title":"To delete"}`)
	http.Post(srv.URL+"/api/tasks", "application/json", body)

	// Delete
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -v -run TestCreateTask`
Expected: FAIL — `RegisterHandlers` not defined

- [ ] **Step 3: Implement handlers**

```go
// handlers.go
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
			handleDelete(w, r, store, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func handleList(w http.ResponseWriter, r *http.Request, store *Store) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.List())
}

type createRequest struct {
	Title string `json:"title"`
}

func handleCreate(w http.ResponseWriter, r *http.Request, store *Store) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
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
	task, ok := store.Update(id, req.Title, req.Status)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func handleDelete(w http.ResponseWriter, r *http.Request, store *Store, id int) {
	if !store.Delete(id) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
```

- [ ] **Step 4: Run all tests**

Run: `go test -v`
Expected: All tests PASS

- [ ] **Step 5: Commit**

```bash
git add handlers.go handlers_test.go
git commit -m "feat: add HTTP handlers for task CRUD"
```

---

### Task 3: Go Backend — Main Entry Point

**Files:**
- Create: `go.mod`
- Create: `main.go`

- [ ] **Step 1: Initialize Go module**

Run: `go mod init kanban`

- [ ] **Step 2: Create main.go**

```go
// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	store := NewStore()
	mux := http.NewServeMux()
	RegisterHandlers(mux, store)

	// Serve frontend static files
	mux.Handle("/", http.FileServer(http.Dir("frontend/dist")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build -o /dev/null .`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add go.mod main.go
git commit -m "feat: add main entry point with static file serving"
```

---

### Task 4: SolidJS Frontend — Scaffold

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/index.html`
- Create: `frontend/src/index.tsx`

- [ ] **Step 1: Scaffold the SolidJS project**

Run:
```bash
cd frontend && npx degit solidjs/templates/ts . && npm install
```

- [ ] **Step 2: Add Vite proxy config for dev**

Replace `frontend/vite.config.ts` with:

```ts
import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    port: 3000,
    proxy: {
      "/api": "http://localhost:8080",
    },
  },
  build: {
    target: "esnext",
  },
});
```

- [ ] **Step 3: Verify it builds**

Run: `cd frontend && npm run build`
Expected: `frontend/dist/` created

- [ ] **Step 4: Commit**

```bash
git add frontend/
git commit -m "feat: scaffold SolidJS frontend with Vite proxy"
```

---

### Task 5: SolidJS Frontend — App Component

**Files:**
- Create: `frontend/src/App.tsx`
- Create: `frontend/src/App.css`

- [ ] **Step 1: Write App.tsx**

```tsx
import { createSignal, createResource, For, Show } from "solid-js";
import "./App.css";

interface Task {
  id: number;
  title: string;
  status: string;
}

const COLUMNS = [
  { key: "todo", label: "Todo" },
  { key: "in-progress", label: "In Progress" },
  { key: "done", label: "Done" },
];

const statusOrder = ["todo", "in-progress", "done"];

async function fetchTasks(): Promise<Task[]> {
  const res = await fetch("/api/tasks");
  return res.json();
}

export default function App() {
  const [tasks, { refetch }] = createResource(fetchTasks);
  const [newTitle, setNewTitle] = createSignal("");

  async function addTask(e: Event) {
    e.preventDefault();
    const title = newTitle().trim();
    if (!title) return;
    await fetch("/api/tasks", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title }),
    });
    setNewTitle("");
    refetch();
  }

  async function moveTask(task: Task, direction: -1 | 1) {
    const idx = statusOrder.indexOf(task.status);
    const newStatus = statusOrder[idx + direction];
    if (!newStatus) return;
    await fetch(`/api/tasks/${task.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title: task.title, status: newStatus }),
    });
    refetch();
  }

  async function deleteTask(id: number) {
    await fetch(`/api/tasks/${id}`, { method: "DELETE" });
    refetch();
  }

  return (
    <div class="app">
      <h1>Kanban Board</h1>
      <form class="add-form" onSubmit={addTask}>
        <input
          type="text"
          placeholder="New task..."
          value={newTitle()}
          onInput={(e) => setNewTitle(e.currentTarget.value)}
        />
        <button type="submit">Add</button>
      </form>
      <div class="board">
        <For each={COLUMNS}>
          {(col) => (
            <div class="column">
              <h2>{col.label}</h2>
              <div class="cards">
                <Show when={tasks()}>
                  <For
                    each={tasks()!.filter((t) => t.status === col.key)}
                  >
                    {(task) => (
                      <div class="card">
                        <span class="card-title">{task.title}</span>
                        <div class="card-actions">
                          <Show when={statusOrder.indexOf(task.status) > 0}>
                            <button onClick={() => moveTask(task, -1)}>←</button>
                          </Show>
                          <Show
                            when={
                              statusOrder.indexOf(task.status) <
                              statusOrder.length - 1
                            }
                          >
                            <button onClick={() => moveTask(task, 1)}>→</button>
                          </Show>
                          <button
                            class="delete"
                            onClick={() => deleteTask(task.id)}
                          >
                            ×
                          </button>
                        </div>
                      </div>
                    )}
                  </For>
                </Show>
              </div>
            </div>
          )}
        </For>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Write App.css**

```css
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background: #f0f2f5;
  color: #1a1a1a;
}

.app {
  max-width: 1100px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

h1 {
  text-align: center;
  margin-bottom: 1.5rem;
  font-size: 1.8rem;
}

.add-form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 2rem;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

.add-form input {
  flex: 1;
  padding: 0.5rem 0.75rem;
  border: 1px solid #ccc;
  border-radius: 6px;
  font-size: 1rem;
}

.add-form button {
  padding: 0.5rem 1.25rem;
  background: #4a90d9;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
}

.add-form button:hover {
  background: #357abd;
}

.board {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.column {
  background: #e4e7eb;
  border-radius: 8px;
  padding: 1rem;
  min-height: 300px;
}

.column h2 {
  font-size: 1rem;
  margin-bottom: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #555;
}

.cards {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.card {
  background: white;
  border-radius: 6px;
  padding: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
}

.card-title {
  flex: 1;
  word-break: break-word;
}

.card-actions {
  display: flex;
  gap: 0.25rem;
  flex-shrink: 0;
}

.card-actions button {
  width: 28px;
  height: 28px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #f8f8f8;
  cursor: pointer;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-actions button:hover {
  background: #e8e8e8;
}

.card-actions button.delete {
  color: #d9534f;
  border-color: #d9534f33;
}

.card-actions button.delete:hover {
  background: #d9534f;
  color: white;
}
```

- [ ] **Step 3: Update index.tsx to use App**

Replace `frontend/src/index.tsx` with:

```tsx
import { render } from "solid-js/web";
import App from "./App";

render(() => <App />, document.getElementById("root")!);
```

- [ ] **Step 4: Build and verify**

Run: `cd frontend && npm run build`
Expected: Builds successfully

- [ ] **Step 5: Commit**

```bash
git add frontend/src/App.tsx frontend/src/App.css frontend/src/index.tsx
git commit -m "feat: add kanban board UI with task management"
```

---

### Task 6: Integration — Build & Run

- [ ] **Step 1: Build frontend**

Run: `cd frontend && npm run build`

- [ ] **Step 2: Run the server**

Run: `go run .`
Expected: "Server running on http://localhost:8080"

- [ ] **Step 3: Verify in browser**

Open http://localhost:8080 — should see the kanban board with three columns.
Add a task, move it between columns, delete it.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "feat: complete kanban app with Go backend and SolidJS frontend"
```
