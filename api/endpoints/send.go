package endpoints

import (
	"whisper-api/communication"

	"github.com/gin-gonic/gin"
)


type MessageRequest struct {
	DeviceID string `json:"device_id"`
	Message string `json:"message"` 
	Subscribers []string `json:"subscribers"` 
}

type SendEndpoint struct {
	Com *communication.Communication
}

func (endpoint SendEndpoint) Handle(c *gin.Context) {
	var request MessageRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.String(400, err.Error())
	}

	conn := communication.Clients[request.DeviceID]
	if conn == nil {
		c.JSON(404, gin.H{
			"error": "device not connected",
		})
		return
	}

	if err := conn.WriteJSON(request); err != nil {
		c.JSON(500, gin.H{
			"error": "failed to send message",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"message": "message sent",
	})
}
