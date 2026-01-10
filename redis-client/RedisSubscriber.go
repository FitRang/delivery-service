package cache

import (
	"context"
	"encoding/json"
	"github.com/FitRang/delivery-service/connections"
	"log"
)

type UserIdentity struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Message struct {
	Sender   UserIdentity `json:"sender"`
	Receiver UserIdentity `json:"receiver"`
	Message  string       `json:"message"`
}

func (rdb *RedisClient) StartRedisSubscriber(
	ctx context.Context,
	hub *connections.Hub,
) {
	sub := rdb.rdb.PSubscribe(ctx, "user:*")
	ch := sub.Channel()

	log.Println("Redis subscriber started (pattern: user:*)")

	for msg := range ch {
		var notif Message
		if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		final, err := finalMessage(notif)
		if err != nil {
			log.Println("invalid message:", err)
			continue
		}

		hub.SendToUser(notif.Receiver.Email, []byte(final))
	}
}
