package api

import (
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort = "8000"
)

func StartServer(wg *sync.WaitGroup) {

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}

	isProduction := os.Getenv("PROD_ENV")
	if isProduction == "" {
		isProduction = "True"
	}

	log.Println("Production env = " + isProduction)

	if isProduction == "True" {
		gin.SetMode(gin.ReleaseMode)
	}

	imgPath = os.Getenv("IMG_PATH")
	if imgPath == "" {
		imgPath = "/var/gobot/img/"
	}

	r := gin.Default()

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/refresh-token", TokenAuthMiddleware(), RefreshToken)
		auth.POST("/logout", TokenAuthMiddleware(), Logout)
	}

	users := api.Group("/users")
	{
		users.POST("/create", TokenAuthMiddleware(), CreateUser)
		users.POST("/update", TokenAuthMiddleware(), UpdateUser)
		users.POST("/delete", TokenAuthMiddleware(), DeleteUser)
		users.GET("/list", TokenAuthMiddleware(), GetAllUsers)
		users.GET("/get", TokenAuthMiddleware(), GetUser)
	}

	messages := api.Group("/messages")
	{
		messages.POST("/create", TokenAuthMiddleware(), CreateBotMessage)
		messages.GET("/list", TokenAuthMiddleware(), GetAllBotMessages)
		messages.POST("/update", TokenAuthMiddleware(), CreateBotMessage)
		messages.POST("/delete", TokenAuthMiddleware(), DeleteMessage)
	}

	images := api.Group("/images")
	{
		images.POST("/create", TokenAuthMiddleware(), CreateImage)
		images.POST("/delete", TokenAuthMiddleware(), DeleteImage)
		images.GET("/list", TokenAuthMiddleware(), GetAllImages)
		images.POST("/upload", TokenAuthMiddleware(), UploadHandler(imgPath))
	}

	requests := api.Group("/requests")
	{
		requests.GET("/list", TokenAuthMiddleware(), GetAllRequests)
		requests.GET("/get", TokenAuthMiddleware(), GetSomeRequests)
	}

	r.GET("/alive", TestAlive)

	log.Fatal(r.Run(":" + port))

	defer wg.Done()
}
