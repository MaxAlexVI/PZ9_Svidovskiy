package task

import (
	"errors"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Repo struct {
	mu   sync.RWMutex
	data map[int64]Task
	next int64
}

func NewRepo() *Repo {
	return &Repo{
		next: 4,
		data: map[int64]Task{
			1: {ID: 1, Title: "Первая задача", Description: "Учебный пример", Done: false},
			2: {ID: 2, Title: "Вторая задача", Description: "Проверка API", Done: true},
			3: {ID: 3, Title: "Изучить Redis", Description: "Кэширование задач", Done: false},
		},
	}
}

func (r *Repo) List() []Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Task, 0, len(r.data))
	for _, t := range r.data {
		result = append(result, t)
	}
	return result
}

func (r *Repo) GetByID(id int64) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.data[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return t, nil
}

func (r *Repo) Create(title string, description string) Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	t := Task{
		ID:          r.next,
		Title:       title,
		Description: description,
		Done:        false,
	}
	r.data[t.ID] = t
	r.next++
	return t
}

func (r *Repo) Update(t Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[t.ID]; !ok {
		return ErrTaskNotFound
	}
	r.data[t.ID] = t
	return nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrTaskNotFound
	}
	delete(r.data, id)
	return nil
}
