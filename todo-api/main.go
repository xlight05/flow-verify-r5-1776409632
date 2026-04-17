package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asdlc/todo-api/internal/handlers"
	"github.com/asdlc/todo-api/internal/middleware"
	"github.com/asdlc/todo-api/internal/store"
)

func main() {
	dataFile := os.Getenv("TODO_DATA_FILE")
	if dataFile == "" {
		dataFile = "/data/todos.json"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	s, err := store.New(dataFile)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	h := handlers.New(s)
	handler := middleware.Recover(middleware.Logging(h.Routes()))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("todo-api listening on :%s (data file: %s)", port, dataFile)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
