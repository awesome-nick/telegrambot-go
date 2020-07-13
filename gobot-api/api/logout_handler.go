package api

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	au, err := ExtractTokenMetadata(c.Request)
	if err != nil {
	   c.JSON(http.StatusUnauthorized, "unauthorized")
		log.Println(err)
	   return
	}
	deleted, delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
	   c.JSON(http.StatusUnauthorized, "unauthorized")
	   return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
  }