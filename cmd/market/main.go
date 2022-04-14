package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"market/internal/infrastructure/exporter"
	"market/internal/infrastructure/postgres"
	"market/internal/interfaces/handlers"
	"market/internal/interfaces/repository"
	"market/internal/usecases/storage"

	"github.com/D3vR4pt0rs/logger"

	"github.com/gorilla/mux"
)

func main() {

	config := postgres.Config{
		Username: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Ip:       os.Getenv("POSTGRES_IP"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DATABASE"),
	}

	psClient := postgres.New(config)
	expClient := exporter.New()

	repo := repository.New(psClient, expClient)

	application := storage.New(repo)

	router := mux.NewRouter()
	handlers.Make(router, application)
	srv := &http.Server{
		Addr:    ":1337",
		Handler: router,
	}

	go func() {
		listener := make(chan os.Signal, 1)
		signal.Notify(listener, os.Interrupt, syscall.SIGTERM)
		fmt.Println("Received a shutdown signal:", <-listener)

		if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			logger.Error.Println("Failed to gracefully shutdown ", err)
		}
	}()

	logger.Info.Println("[*]  Listening...")
	if err := srv.ListenAndServe(); err != nil {
		logger.Error.Println("Failed to listen and serve ", err)
	}

	logger.Critical.Println("Server shutdown")
}
