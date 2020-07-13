package api

import (
	"gobot/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var GetAllRequests = func(c *gin.Context) {
	requests := database.GetAllRequests()
	c.JSON(http.StatusOK, requests)
}

var GetSomeRequests = func(c *gin.Context) {

	f, err := strconv.ParseUint(c.Query("from"), 10, 64)
	if err != nil {
		f = 0
	}

	t, err := strconv.ParseUint(c.Query("to"), 10, 64)
	if err != nil {
		t = 0
	}

	if t == 0 {
		if f == 0 { // if so, we'll provide all requests
			requests := database.GetAllRequests()
			c.JSON(http.StatusOK, requests)
			return
		} else {
			requests := database.GetAllRequests(f)
			c.JSON(http.StatusOK, requests)
			return
		}
	}
	requests := database.GetAllRequests(f, t)
	c.JSON(http.StatusOK, requests)
}
