package database

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *gorm.DB
)

func init() {

	// e := godotenv.Load(os.ExpandEnv("$GOPATH/src/gobot/.env"))
	// if e != nil {
	// 	fmt.Print(e)
	// }

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbType := os.Getenv("DB_TYPE")
	dbPort := os.Getenv("DB_PORT")
	adminPass := os.Getenv("ADMIN_PASS")

	dsn := url.URL{
		User:     url.UserPassword(username, password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:     dbName,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	conn, err := gorm.Open(dbType, dsn.String())

	if err != nil {
		log.Fatalln(err)
	}

	db = conn
	db.Debug().AutoMigrate(&BotMessage{}, &Request{}, &APIUser{}, &Image{})
	CreateAdmin(adminPass)

}

func GetDB() *gorm.DB {
	return db
}
