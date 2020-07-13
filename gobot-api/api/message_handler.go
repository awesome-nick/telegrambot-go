package api

import (
	"gobot/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

var CreateBotMessage = func(c *gin.Context) {
	var b *database.BotMessage
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	b.Create()

	c.JSON(http.StatusCreated, b)
}

var GetAllBotMessages = func(c *gin.Context) {
	messages := database.GetBotMessages()
	c.JSON(http.StatusOK, messages)
}

var DeleteMessage = func(c *gin.Context) {
	var b *database.BotMessage
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	if ok := b.Delete(); ok {
		c.JSON(http.StatusOK, b)
	} else {
		c.JSON(http.StatusAccepted, "Message "+b.MessageId+" not found")
	}

}
