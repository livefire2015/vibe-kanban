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
                  <For each={tasks()!.filter((t) => t.status === col.key)}>
                    {(task) => (
                      <div class="card">
                        <span class="card-title">{task.title}</span>
                        <div class="card-actions">
                          <Show when={statusOrder.indexOf(task.status) > 0}>
                            <button onClick={() => moveTask(task, -1)}>
                              &#8592;
                            </button>
                          </Show>
                          <Show
                            when={
                              statusOrder.indexOf(task.status) <
                              statusOrder.length - 1
                            }
                          >
                            <button onClick={() => moveTask(task, 1)}>
                              &#8594;
                            </button>
                          </Show>
                          <button
                            class="delete"
                            onClick={() => deleteTask(task.id)}
                          >
                            &times;
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
