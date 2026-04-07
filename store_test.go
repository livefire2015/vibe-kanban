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
