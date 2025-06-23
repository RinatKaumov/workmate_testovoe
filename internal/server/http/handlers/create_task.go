package handlers

import (
	"awesomeProject36/internal/domain/service"
	"encoding/json"
	"log"
	"net/http"
)

func HandleCreateTask(taskService *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Description string `json:"description"`
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("ошибка декодирования JSON: %v", err)
			http.Error(w, "неверный формат JSON", http.StatusBadRequest)
			return
		}

		// Валидация описания задачи
		if err := validateTaskDescription(req.Description); err != nil {
			log.Printf("ошибка валидации описания: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task, err := taskService.CreateTask(r.Context(), req.Description)
		if err != nil {
			log.Printf("ошибка создания задачи: %v", err)
			http.Error(w, "не удалось создать задачу", http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusCreated, task)
	}
}
