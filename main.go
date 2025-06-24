package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	repo "github.com/RinatKaumov/workmate_testovoe/internal/dal/repository"
	"github.com/RinatKaumov/workmate_testovoe/internal/domain/service"
	server "github.com/RinatKaumov/workmate_testovoe/internal/server/http"
)

func main() {

	repository := repo.NewTaskRepositoryInMemory()
	taskService := service.NewTaskService(repository)
	srv := server.NewServer(taskService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort("", port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	// Канал для сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Сервер запущен на %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	<-quit
	log.Println("Остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	// Ждем завершения всех фоновых задач
	log.Println("Ожидание завершения фоновых задач...")
	taskService.WaitForTasks()

	log.Println("Сервер остановлен корректно")
}
