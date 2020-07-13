package api

import (
	"gobot/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

var CreateImage = func(c *gin.Context) {

	var i database.Image
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	i.Create()
	c.JSON(http.StatusCreated, i)
}

var DeleteImage = func(c *gin.Context) {

	var i database.Image
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if d := i.Delete(); d {
		c.JSON(http.StatusOK, i)
	} else {
		c.JSON(http.StatusAccepted, "Image "+i.Command+" not found")
	}

}

var GetAllImages = func(c *gin.Context) {
	if is := database.GetAllImages(); is != nil {
		c.JSON(http.StatusOK, is)
	} else {
		c.JSON(http.StatusInternalServerError, "Error while retrieving images")
	}

}
