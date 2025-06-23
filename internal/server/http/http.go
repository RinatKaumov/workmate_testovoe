package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"awesomeProject36/internal/domain/service"
	handlers "awesomeProject36/internal/server/http/handlers"
)

// Server обёртка для сервиса задач и роутера
type Server struct {
	service *service.TaskService
	router  *chi.Mux
}

func NewServer(taskService *service.TaskService) *Server {
	s := &Server{
		service: taskService,
		router:  chi.NewRouter(),
	}
	s.configureRoutes()
	return s
}

func (s *Server) configureRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Get("/tasks", handlers.HandleListTasks(s.service))
	s.router.Post("/tasks", handlers.HandleCreateTask(s.service))
	s.router.Get("/tasks/{id}", handlers.HandleGetTask(s.service))
	s.router.Delete("/tasks/{id}", handlers.HandleDeleteTask(s.service))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
