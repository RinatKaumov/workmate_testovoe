package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/RinatKaumov/workmate_testovoe/internal/domain/service"
	handler "github.com/RinatKaumov/workmate_testovoe/internal/server/http/handler"
)

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
	s.router.Get("/tasks", handler.ListTasks(s.service))
	s.router.Post("/tasks", handler.CreateTask(s.service))
	s.router.Get("/tasks/{id}", handler.GetTask(s.service))
	s.router.Delete("/tasks/{id}", handler.DeleteTask(s.service))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
