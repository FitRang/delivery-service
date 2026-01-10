package cache

import "encoding/json"

type FinalMessage struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Message  string `json:"message"`
}

func finalMessage(notif Message) ([]byte, error) {
	message := FinalMessage{
		Receiver: notif.Receiver.Username,
		Sender:   notif.Sender.Username,
		Message:  notif.Message,
	}
	return json.Marshal(message)
}
