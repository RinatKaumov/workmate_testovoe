package handlers

import (
	"errors"
	"strings"
)

// validateTaskDescription проверяет корректность описания задачи
func validateTaskDescription(description string) error {
	// Проверка на пустую строку
	if strings.TrimSpace(description) == "" {
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

// validateTaskID проверяет корректность ID задачи
func validateTaskID(id string) error {
	// Проверка на пустой ID
	if strings.TrimSpace(id) == "" {
		return errors.New("ID задачи не может быть пустым")
	}

	// Проверка минимальной длины (UUID обычно 36 символов)
	if len(id) < 10 {
		return errors.New("ID задачи слишком короткий")
	}

	// Проверка максимальной длины
	if len(id) > 100 {
		return errors.New("ID задачи слишком длинный")
	}

	return nil
}
