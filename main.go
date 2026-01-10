package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FitRang/delivery-service/connections"
	cache "github.com/FitRang/delivery-service/redis-client"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := connections.NewHub()

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		log.Fatal("REDIS_URI not set")
	}

	rdb, err := cache.NewRedisClient(redisURI)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Redis subscriber goroutine")
	go rdb.StartRedisSubscriber(ctx, hub)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connections.ServeWS(hub, w, r)
	})

	server := &http.Server{
		Addr:    ":8084",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("WebSocket server running on :8084")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	hub.Close()

	if err := rdb.Close(); err != nil {
		log.Printf("Redis close error: %v", err)
	}

	log.Println("Delivery service shut down gracefully")
}
