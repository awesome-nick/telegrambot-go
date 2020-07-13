package api

import (
	"encoding/base64"
	"errors"
	"gobot/database"
	"gobot/redis"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

var apiSecret, apiRefreshSecret string

func init() {
	// if err := godotenv.Load(os.ExpandEnv("$GOPATH/src/gobot/.env")); err != nil {
	// 	log.Println("Error loading .env file! Default API_SECRET & API_REFRESH_SECRET will be used!")
	// }

	apiSecret = os.Getenv("API_SECRET")
	if apiSecret == "" {
		apiSecret = "dummyapiseretkey"
	}

	apiRefreshSecret = os.Getenv("API_REFRESH_SECRET")
	if apiRefreshSecret == "" {
		apiRefreshSecret = "ijufsrcijriuiuwirouiosriwrci"
	}
}

func getBasicAuth(s string) ([]string, error) {

	a := strings.SplitN(s, " ", 2)
	if len(a) != 2 || a[0] != "Basic" {
		return []string{}, errors.New("Not auth string")
	}

	payload, err := base64.StdEncoding.DecodeString(a[1])
	if err != nil {
		return []string{}, err
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if pair[0] == "" || pair[1] == "" {
		return []string{}, errors.New("Empty username or pass")
	}

	return pair, nil

}

var Login = func(c *gin.Context) {

	pair, err := getBasicAuth(c.GetHeader("Authorization"))

	if err != nil {
		c.JSON(http.StatusUnauthorized, "authorization failed")
		return
	}

	apiUser := database.APIUser{
		Username: pair[0],
		Password: pair[1],
	}

	aUser := database.GetAPIUserByUserName(apiUser.Username)

	if aUser == nil {
		c.JSON(http.StatusUnauthorized, "User with ID "+pair[0]+" not found in database")
		return
	}
	passwdOk := database.CompareHashAndPasswd(aUser.Password, apiUser.Password)

	//compare the user from the request, with the one we defined:
	if aUser.Username != apiUser.Username || !passwdOk {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(aUser.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := CreateAuth(aUser.ID, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

var CreateToken = func(userid uint64) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// Creating Access Token
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("API_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	redisClient := redis.GetRedis()

	errAccess := redisClient.Set(redis.Context, td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := redisClient.Set(redis.Context, td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
