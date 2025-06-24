package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          uuid.UUID
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
	Result      string
	Error       string
}
