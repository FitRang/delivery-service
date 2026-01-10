package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/FitRang/delivery-service/connections"
	cache "github.com/FitRang/delivery-service/redis-client"
)

func main() {
	ctx := context.Background()
	hub := connections.NewHub()

	rdb, err := cache.NewRedisClient(os.Getenv("REDIS_URI"))
	if err != nil {
		log.Fatal(err)
	}

	go rdb.StartRedisSubscriber(ctx, hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connections.ServeWS(hub, w, r)
	})

	log.Println("WebSocket server running on :8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
