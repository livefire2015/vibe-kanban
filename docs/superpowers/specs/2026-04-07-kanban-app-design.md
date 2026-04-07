# Kanban App Design

## Overview
Simple kanban board with SolidJS frontend and Go backend. In-memory storage, three fixed columns.

## Backend (Go, standard library)
- `net/http` router, no external dependencies
- In-memory store: `[]Task` protected by `sync.Mutex`
- Serves frontend static files from `frontend/dist/`

### Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/tasks | List all tasks |
| POST | /api/tasks | Create task (body: `{title}`, defaults status to `todo`) |
| PUT | /api/tasks/{id} | Update task (body: `{title?, status?}`) |
| DELETE | /api/tasks/{id} | Delete task |

### Data Model
```go
type Task struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status string `json:"status"` // "todo" | "in-progress" | "done"
}
```

## Frontend (SolidJS + Vite)
- Three columns: Todo, In Progress, Done
- Each card: title, move-left/move-right buttons, delete button
- Add-task form at the top
- Minimal CSS, no component library, no drag-and-drop

## Non-goals
- No auth, no routing, no database, no drag-and-drop
