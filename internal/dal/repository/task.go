package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/RinatKaumov/workmate_testovoe/internal/dal/entity"
)

var ErrNotFound = errors.New("not found")

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
}

type TaskRepositoryInMemory struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*entity.Task
}

func NewTaskRepositoryInMemory() *TaskRepositoryInMemory {
	return &TaskRepositoryInMemory{
		tasks: make(map[uuid.UUID]*entity.Task, 100),
	}
}

func (r *TaskRepositoryInMemory) Create(ctx context.Context, task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepositoryInMemory) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return task, nil
}

func (r *TaskRepositoryInMemory) DeleteByID(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return ErrNotFound
	}

	delete(r.tasks, id)
	return nil
}

func (r *TaskRepositoryInMemory) List(ctx context.Context) ([]*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*entity.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		result = append(result, task)
	}
	return result, nil
}

func (r *TaskRepositoryInMemory) Update(ctx context.Context, task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}
