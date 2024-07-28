package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weldonkipchirchir/rental_listing/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	stopChan := make(chan os.Signal, 1) // Buffered channel to fix SA1017
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start("0.0.0.0:8000"); err != nil {
			log.Fatal("Cannot start server:", err)
		}
	}()

	<-stopChan

	log.Printf("Shutting down server...")

	// Create a context with a timeout to allow the server to shut down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shut down:", err)
	}

	log.Println("Server exiting")
}
