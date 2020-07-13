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
	GetDB().Create(request)
}

func GetRequestsByChatId(id uint64) []*Request {

	r := make([]*Request, 0)
	err := GetDB().Table("requests").Where("chat_id = ?", id).Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return r
}

func GetAllRequests(p ...uint64) []*Request {

	r := make([]*Request, 0)
	switch len(p) {
	case 1:
		err := GetDB().Table("requests").Find(&r, "id >= ?", p[0]).Error
		if err != nil {
			return nil
		}
		return r
	case 2:
		err := GetDB().Table("requests").Find(&r, "id >= ? AND id <= ?", p[0], p[1]).Error
		if err != nil {
			return nil
		}
		return r
	default:
		err := GetDB().Table("requests").Find(&r).Error
		if err != nil {
			return nil
		}
		return r
	}
}
