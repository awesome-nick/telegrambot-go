package api

import (
	"gobot/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

var TestAlive = func(c *gin.Context) {
	t := database.GetAPIUsers()
	if len(t) != 0 {
		c.JSON(http.StatusOK, "alive")
	} else {
		c.JSON(http.StatusInternalServerError, "Oops")
	}
}
