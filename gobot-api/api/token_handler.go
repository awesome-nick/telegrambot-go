package api

import (
	"fmt"
	"gobot/redis"
	"net/http"
	"os"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

func init() {

	apiSecret = os.Getenv("API_SECRET")
	if apiSecret == "" {
		apiSecret = "dummyapisecretkey"
	}

	apiRefreshSecret = os.Getenv("API_REFRESH_SECRET")
	if apiRefreshSecret == "" {
		apiRefreshSecret = "ijufsrcijriuiuwirouiosriwrci"
	}
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tString := ExtractToken(r)
	t, err := jwt.Parse(tString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(apiSecret), nil
	})
	if err != nil {

		return nil, err
	}
	return t, nil
}

func TokenValid(r *http.Request) error {
	t, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := t.Claims.(jwt.Claims); !ok && !t.Valid {
		return err
	}
	return nil
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := redis.GetRedis().Get(redis.Context, authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func DeleteAuth(givenUuid string) (int64, error) {
	d, err := redis.GetRedis().Del(redis.Context, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return d, nil
}
