package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"

	"github.com/RinatKaumov/workmate_testovoe/internal/dal/entity"
	"github.com/RinatKaumov/workmate_testovoe/internal/dal/repository"
	"github.com/RinatKaumov/workmate_testovoe/internal/domain/model"
)

var ErrTaskNotFound = errors.New("task not found")

// taskRepository — это интерфейс, описывающий, что именно нужно TaskService
// от хранилища задач.
type taskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
}

// TaskService инкапсулирует бизнес-логику задач.
type TaskService struct {
	repo taskRepository
	wg   sync.WaitGroup // для отслеживания активных фоновых задач
}

// NewTaskService создает новый сервис задач.
func NewTaskService(repo taskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// WaitForTasks ждет завершения всех активных фоновых задач
func (s *TaskService) WaitForTasks() {
	s.wg.Wait()
}

// CreateTask создает задачу и запускает ее асинхронно.
func (s *TaskService) CreateTask(ctx context.Context, description string) (*model.Task, error) {
	now := time.Now().UTC()
	entityTask := &entity.Task{
		ID:          uuid.New(),
		Description: description,
		Status:      entity.StatusPending,
		CreatedAt:   now,
	}

	err := s.repo.Create(ctx, entityTask)
	if err != nil {
		return nil, err
	}

	// Запускаем задачу в фоне с отслеживанием
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runTask(entityTask.ID)
	}()

	return toModel(entityTask)
}

// runTask — имитация долгой задачи (3-5 минут)
func (s *TaskService) runTask(id uuid.UUID) {
	ctx := context.Background()
	// Обновляем статус задачи на running
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("ошибка получения задачи в runTask: %v\n", err)
		return
	}

	now := time.Now().UTC()
	task.Status = entity.StatusRunning
	task.StartedAt = &now
	if err := s.repo.Update(ctx, task); err != nil {
		fmt.Printf("ошибка обновления задачи (running): %v\n", err)
		return
	}

	// Долгая работа (например, 3-5 минут)
	duration := 3*time.Minute + time.Duration(rand.Intn(3))*time.Minute
	time.Sleep(duration)

	// По завершении обновляем статус и результат у уже имеющегося указателя
	finished := time.Now().UTC()
	task.Status = entity.StatusCompleted
	task.FinishedAt = &finished
	task.Result = "result data here"

	if err := s.repo.Update(ctx, task); err != nil {
		fmt.Printf("ошибка обновления задачи (completed): %v\n", err)
		return
	}
}

// GetTask возвращает задачу по ID.
func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	entityTask, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgErrors.WithStack(ErrTaskNotFound)
		}
		return nil, err // Пробрасываем другие ошибки как есть
	}
	return toModel(entityTask)
}

// DeleteTask удаляет задачу по ID.
func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgErrors.WithStack(ErrTaskNotFound)
		}
		return err
	}
	return nil
}

// ListTasks возвращает список всех задач.
func (s *TaskService) ListTasks(ctx context.Context) ([]*model.Task, error) {
	entityTasks, err := s.repo.List(ctx)
	if err != nil {
		return nil, err // Пробрасываем ошибку репозитория наверх
	}

	if len(entityTasks) == 0 {
		return []*model.Task{}, nil
	}

	tasks := make([]*model.Task, 0, len(entityTasks))
	for _, taskEntity := range entityTasks {
		taskModel, err := toModel(taskEntity)
		if err != nil {
			return nil, err // Прерываем выполнение и возвращаем ошибку
		}
		tasks = append(tasks, taskModel)
	}
	return tasks, nil
}

// toModel конвертирует entity.Task в model.Task
func toModel(e *entity.Task) (*model.Task, error) {
	if e == nil {
		return nil, errors.New("cannot convert nil entity task to model")
	}

	var duration *int64
	if e.StartedAt != nil {
		var d int64
		if e.FinishedAt != nil {
			d = int64(e.FinishedAt.Sub(*e.StartedAt).Seconds())
		} else {
			now := time.Now().UTC()
			d = int64(now.Sub(*e.StartedAt).Seconds())
		}
		duration = &d
	}

	return &model.Task{
		ID:          e.ID,
		Description: e.Description,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		StartedAt:   e.StartedAt,
		FinishedAt:  e.FinishedAt,
		Result:      e.Result,
		Error:       e.Error,
		Duration:    duration,
	}, nil
}

// toEntity — если понадобится обратное преобразование
func toEntity(m *model.Task) (*entity.Task, error) {
	if m == nil {
		return nil, errors.New("cannot convert nil model task to entity")
	}
	return &entity.Task{
		ID:          m.ID,
		Description: m.Description,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		StartedAt:   m.StartedAt,
		FinishedAt:  m.FinishedAt,
		Result:      m.Result,
		Error:       m.Error,
	}, nil
}
