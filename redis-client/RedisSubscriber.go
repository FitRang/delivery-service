package cache

import (
	"context"
	"encoding/json"
	"log"
	"github.com/FitRang/delivery-service/connections"
)

type NotificationMessage struct {
	UserID  string          `json:"user_id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (rdb *RedisClient) StartRedisSubscriber(ctx context.Context, hub *connections.Hub) {
	sub := rdb.rdb.Subscribe(ctx, "users:*")
	ch := sub.Channel()

	for msg := range ch {
		var notif NotificationMessage
		if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		hub.SendToUser(notif.UserID, []byte(msg.Payload))
	}
}
