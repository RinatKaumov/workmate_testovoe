package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"awesomeProject36/internal/dal/entity"
	"awesomeProject36/internal/dal/repository"
	"awesomeProject36/internal/domain/model"
)

// TaskService инкапсулирует бизнес-логику задач.
type TaskService struct {
	repo repository.TaskRepository
}

// NewTaskService создает новый сервис задач.
func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask создает задачу и запускает ее асинхронно.
func (s *TaskService) CreateTask(ctx context.Context, description string) (*model.Task, error) {
	// Валидация на уровне бизнес-логики
	if err := s.validateTaskCreation(description); err != nil {
		return nil, err
	}

	entityTask, err := s.repo.Create(description)
	if err != nil {
		return nil, err
	}

	// Запускаем задачу в фоне
	go s.runTask(entityTask.ID)

	return toModel(entityTask), nil
}

// validateTaskCreation проверяет возможность создания задачи
func (s *TaskService) validateTaskCreation(description string) error {
	// Дополнительная бизнес-логика валидации
	// Например, проверка на дубликаты описаний
	// Или проверка лимитов на количество задач

	// Здесь можно добавить более сложную логику
	// Например, проверку на максимальное количество активных задач
	// или проверку на дубликаты описаний

	return nil
}

// runTask — имитация долгой задачи (3-5 минут)
func (s *TaskService) runTask(id string) {
	// Обновляем статус задачи на running
	task, ok := s.repo.GetByID(id)
	if !ok {
		return // задача пропала из репо
	}

	now := time.Now().UTC()
	task.Status = entity.StatusRunning
	task.StartedAt = &now
	s.repo.Update(task)

	// Долгая работа (например, 3-5 минут)
	duration := 3*time.Minute + time.Duration(rand.Intn(3))*time.Minute
	time.Sleep(duration)

	// По завершении обновляем статус и результат
	task, ok = s.repo.GetByID(id)
	if !ok {
		return
	}

	finished := time.Now().UTC()
	task.Status = entity.StatusCompleted
	task.FinishedAt = &finished
	result := "result data here"
	task.Result = &result

	err := s.repo.Update(task)
	if err != nil {
		return
	}
}

// GetTask возвращает задачу по ID.
func (s *TaskService) GetTask(ctx context.Context, id string) (*model.Task, error) {
	entityTask, ok := s.repo.GetByID(id)
	if !ok {
		return nil, errors.New("task not found")
	}
	return toModel(entityTask), nil
}

// DeleteTask удаляет задачу по ID.
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	deleted := s.repo.DeleteByID(id)
	if !deleted {
		return errors.New("task not found")
	}
	return nil
}

// ListTasks возвращает список всех задач.
func (s *TaskService) ListTasks(ctx context.Context) ([]*model.Task, error) {
	entityTasks := s.repo.List()
	tasks := make([]*model.Task, len(entityTasks))
	for i, task := range entityTasks {
		tasks[i] = toModel(task)
	}
	return tasks, nil
}

// toModel конвертирует entity.Task в model.Task
func toModel(e *entity.Task) *model.Task {
	if e == nil {
		return nil
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
		Status:      model.TaskStatus(e.Status),
		CreatedAt:   e.CreatedAt,
		StartedAt:   e.StartedAt,
		FinishedAt:  e.FinishedAt,
		Result:      e.Result,
		Error:       e.Error,
		Duration:    duration,
	}
}

// toEntity — если понадобится обратное преобразование
func toEntity(m *model.Task) *entity.Task {
	if m == nil {
		return nil
	}
	return &entity.Task{
		ID:          m.ID,
		Description: m.Description,
		Status:      entity.TaskStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		StartedAt:   m.StartedAt,
		FinishedAt:  m.FinishedAt,
		Result:      m.Result,
		Error:       m.Error,
	}
}
