package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tabrej2001/student-api/internal/config"
	"github.com/tabrej2001/student-api/internal/http/handlers/students"
	"github.com/tabrej2001/student-api/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env:", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New(storage))
	router.HandleFunc("GET /api/students/{id}", students.GetStudentByID(storage))

	fmt.Println("Server started")
	slog.Info("Server started at por t", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start the server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
