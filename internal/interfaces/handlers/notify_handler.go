package handlers

import (

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

type NotifyHandler struct{
	client *linebot.Client
}

func NewNotifyHandler(lineClient *linebot.Client) *NotifyHandler {
	return	&NotifyHandler{client: lineClient}
}

func (h *NotifyHandler) SendMessage(c *gin.Context){
	message := linebot.NewTextMessage("Hello World")
	_, err := h.client.BroadcastMessage(message).Do()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error", 
			"message": "Failed to broadcast message",
		})
	}

	c.JSON(200, gin.H{
		"status": "success",
		"message": "Broadcast message",
	})
}
