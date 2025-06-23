package entity

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
	Result      *string
	Error       *string
}

func Now() time.Time {
	return time.Now().UTC()
}
