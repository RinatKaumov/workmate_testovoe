package model

import (
	"time"

	"github.com/RinatKaumov/workmate_testovoe/internal/dal/entity"
	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID         `json:"id"`
	Description string            `json:"description"`
	Status      entity.TaskStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	FinishedAt  *time.Time        `json:"finished_at,omitempty"`
	Result      string            `json:"result"`
	Error       string            `json:"error"`
	Duration    *int64            `json:"duration,omitempty"`
}
