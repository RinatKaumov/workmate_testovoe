package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/RinatKaumov/workmate_testovoe/internal/domain/service"

	"github.com/go-chi/chi/v5"
)

func GetTask(taskService *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := validateTaskID(idStr)
		if err != nil {
			log.Printf("ошибка валидации ID: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task, err := taskService.GetTask(r.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrTaskNotFound) {
				log.Printf("ошибка получения задачи: %v", err)
				http.Error(w, "задача не найдена", http.StatusNotFound)
				return
			}
			log.Printf("неожиданная ошибка получения задачи: %v", err)
			http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
		WriteJSON(w, http.StatusOK, task)
	}
}
