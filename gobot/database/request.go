package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Request struct {
	gorm.Model
	MessageId   uint64 `json:"message_id"`
	UserId      uint64 `json:"user_id"`
	FirstName   string `json:"first_name"`
	Language    string `json:"lang"`
	ChatId      uint64 `json:"chat_id"`
	MessageText string `json:"text"`
	Date        uint64 `json:"timestamp"`
}

func (request *Request) Create() {
	// map[string]interface{} {

	GetDB().Create(request)

	// resp := u.Message(true, "success")
	// resp["income"] = rec
	// return resp
}

func GetRequestsByChatId(id uint64) []*Request {

	records := make([]*Request, 0)
	err := GetDB().Table("requests").Where("chat_id = ?", id).Find(&records).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return records
}

func GetAllRequests() []*Request {

	records := make([]*Request, 0)
	err := GetDB().Table("requests").Find(&records).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return records
}
