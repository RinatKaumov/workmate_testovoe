package handler

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// validateTaskDescription проверяет корректность описания задачи
func validateTaskDescription(description string) error {
	// Проверка на пустую строку
	if description == "" {
		return errors.New("описание задачи не может быть пустым")
	}

	// Проверка минимальной длины
	if len(description) < 3 {
		return errors.New("описание задачи должно содержать минимум 3 символа")
	}

	// Проверка максимальной длины
	if len(description) > 1000 {
		return errors.New("описание задачи слишком длинное (максимум 1000 символов)")
	}

	// Проверка на запрещенные символы (базовая защита от XSS)
	if strings.Contains(strings.ToLower(description), "<script>") {
		return errors.New("описание содержит запрещенные символы")
	}

	return nil
}

func validateTaskID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("ID задачи должен быть корректным UUID")
	}
	return parsedID, nil
}
