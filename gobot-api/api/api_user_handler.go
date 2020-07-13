package api

import (
	"gobot/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

var CreateUser = func(c *gin.Context) {

	var u database.APIUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if isUserCreated := u.Create(); isUserCreated {
		u.Password = "hidden"
		c.JSON(http.StatusCreated, u)
	} else {
		c.JSON(http.StatusOK, "User already exists")
	}

}

var UpdateUser = func(c *gin.Context) {

	var u database.APIUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if found := u.Update(); found {
		c.JSON(http.StatusOK, u)
	} else {
		c.JSON(http.StatusNotFound, "User "+u.Username+" not found")
	}

}

var DeleteUser = func(c *gin.Context) {

	var u database.APIUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if d := u.Delete(); d {
		c.JSON(http.StatusOK, u)
	} else {
		c.JSON(http.StatusAccepted, "User "+u.Username+" not found")
	}

}

var GetAllUsers = func(c *gin.Context) {
	c.JSON(http.StatusOK, database.GetAPIUsers())
}

var GetUser = func(c *gin.Context) {
	var u database.APIUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if user := database.GetAPIUserByUserName(u.Username); user != nil {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusNotFound, "User "+u.Username+" not found")
	}
}
