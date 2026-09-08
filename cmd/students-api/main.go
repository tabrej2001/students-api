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
) 

func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	fmt.Println("Server started")
	slog.Info("Server started at por t", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr: cfg.Address,
		Handler: router,
	}

	go func ()  {
    
	err := server.ListenAndServe()

		if err != nil {
		log.Fatal("Failed to start the server")
	}	
	}()


	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}

    slog.Info("Server shutdown successfully") 
}
