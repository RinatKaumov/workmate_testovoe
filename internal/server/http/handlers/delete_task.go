package handlers

import (
	"awesomeProject36/internal/domain/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HandleDeleteTask(taskService *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		// Валидация ID задачи
		if err := validateTaskID(id); err != nil {
			log.Printf("ошибка валидации ID: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := taskService.DeleteTask(r.Context(), id)
		if err != nil {
			log.Printf("ошибка удаления задачи: %v", err)
			http.Error(w, "задача не найдена", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
