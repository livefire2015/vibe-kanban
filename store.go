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
