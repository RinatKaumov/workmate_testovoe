package handlers

import (
	"log"
	"net/http"

	"github.com/RinatKaumov/workmate_testovoe/internal/domain/service"
)

func HandleListTasks(taskService *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := taskService.ListTasks(r.Context())
		if err != nil {
			log.Printf("ошибка получения списка задач: %v", err)
			http.Error(w, "не удалось получить список задач", http.StatusInternalServerError)
			return
		}
		WriteJSON(w, http.StatusOK, tasks)
	}
}
