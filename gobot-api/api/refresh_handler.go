package api

import (
	"net/http"
	"strconv"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var RefreshToken = func(c *gin.Context) {
	mapToken := map[string]string{}

	if err := c.ShouldBindJSON(&mapToken); err != nil {
	   c.JSON(http.StatusUnprocessableEntity, err.Error())
	   return
	}

	refreshToken := mapToken["refresh_token"]
   
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(apiRefreshSecret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
	   c.JSON(http.StatusUnauthorized, "Refresh token expired")
	   return
	}
	
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
	   c.JSON(http.StatusUnauthorized, err)
	   return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
	   refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
	   if !ok {
		  c.JSON(http.StatusUnprocessableEntity, err)
		  return
	   }
	   userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	   if err != nil {
		  c.JSON(http.StatusUnprocessableEntity, "Error occurred")
		  return
	   }
	   //Delete the previous Refresh Token
	   deleted, delErr := DeleteAuth(refreshUuid)
	   if delErr != nil || deleted == 0 { //if any goes wrong
		  c.JSON(http.StatusUnauthorized, "unauthorized")
		  return
	   }
	  //Create new pairs of refresh and access tokens
	   token, createErr := CreateToken(userId)
	   if  createErr != nil {
		  c.JSON(http.StatusForbidden, createErr.Error())
		  return
	   }
	  //save the tokens metadata to redis
  saveErr := CreateAuth(userId, token)
   if saveErr != nil {
		  c.JSON(http.StatusForbidden, saveErr.Error())
		 return
  }
   tokens := map[string]string{
		 "access_token":  token.AccessToken,
	"refresh_token": token.RefreshToken,
  }
	   c.JSON(http.StatusCreated, tokens)
	} else {
	   c.JSON(http.StatusUnauthorized, "refresh expired")
	}
  }