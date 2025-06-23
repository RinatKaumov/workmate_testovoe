package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	repo "awesomeProject36/internal/dal/repository"
	"awesomeProject36/internal/domain/service"
	server "awesomeProject36/internal/server/http"
)

func main() {

	repository := repo.NewTaskRepositoryInMemory()
	taskService := service.NewTaskService(repository)
	srv := server.NewServer(taskService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

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
	log.Println("Сервер остановлен корректно")
}
