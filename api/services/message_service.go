package services

import (
	"fmt"
	"whisper-api/communication"
)

type MessageRequest struct {
	DeviceID    string   `json:"device_id"`
	Message     string   `json:"message"`
	Subscribers []string `json:"subscribers"`
}

func SendMessage(request MessageRequest) error {
	conn := communication.Clients[request.DeviceID]
	if conn == nil {
		return fmt.Errorf("device not connected")
	}

	if err := conn.WriteJSON(request); err != nil {
		return err
	}

	return nil
}
