package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var Context context.Context = context.Background()

func init() {
	// if err := godotenv.Load(os.ExpandEnv("$GOPATH/src/gobot/.env")); err != nil {
	// 	log.Println("Error loading .env file! Default HTTP port 8000 will be used!")
	// }

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1" //localhost
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPass := os.Getenv("REDIS_PASS")
	if redisPass == "" {
		redisPass = ""
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPass, // no password set
		DB:       0,         // use default DB
	})

	_, err := redisClient.Ping(Context).Result()
	if err != nil {
		panic(err)
	}
}

func GetRedis() *redis.Client {
	return redisClient
}
