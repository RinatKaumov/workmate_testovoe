package repository

import (
	"sync"

	"awesomeProject36/internal/dal/entity"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(description string) (*entity.Task, error)
	GetByID(id string) (*entity.Task, bool)
	DeleteByID(id string) bool
	List() []*entity.Task
	Update(task *entity.Task) error
}

type TaskRepositoryInMemory struct {
	mu    sync.RWMutex
	tasks map[string]*entity.Task
}

func NewTaskRepositoryInMemory() *TaskRepositoryInMemory {
	return &TaskRepositoryInMemory{
		tasks: make(map[string]*entity.Task),
	}
}

func (r *TaskRepositoryInMemory) Create(description string) (*entity.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	task := &entity.Task{
		ID:          id,
		Description: description,
		Status:      entity.StatusPending,
		CreatedAt:   entity.Now(),
	}
	r.tasks[id] = task
	return task, nil
}

func (r *TaskRepositoryInMemory) GetByID(id string) (*entity.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	return task, ok
}

func (r *TaskRepositoryInMemory) DeleteByID(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; ok {
		delete(r.tasks, id)
		return true
	}
	return false
}

func (r *TaskRepositoryInMemory) List() []*entity.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*entity.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		result = append(result, task)
	}
	return result
}

func (r *TaskRepositoryInMemory) Update(task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}
