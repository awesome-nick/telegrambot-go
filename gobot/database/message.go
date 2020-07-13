package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type BotMessage struct {
	gorm.Model
	MessageId   string `json:"message_id" gorm:"unique;not null"`
	MessageText string `json:"text"`
	Command     string `json:"command"`
}

func (m *BotMessage) Create() {
	var mh BotMessage
	if err := GetDB().Where("message_id = ?", m.MessageId).First(&mh).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			GetDB().Create(&m)
		}
	} else {
		GetDB().Model(&mh).Where("message_id = ?", m.MessageId).Update(&m)
	}
}

func (m *BotMessage) Delete() bool {

	var bm BotMessage
	if err := GetDB().Where("message_id = ?", m.MessageId).First(&bm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
	}
	GetDB().Unscoped().Where("message_id = ?", m.MessageId).Delete(&m)
	return true
}

func GetBotMessageByCommand(command string) *BotMessage {

	r := make([]*BotMessage, 0)
	err := GetDB().Table("bot_messages").Where("command = ?", command).Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(r) > 0 {
		return r[0]
	}
	return nil
}

func GetBotMessages() []*BotMessage {

	r := make([]*BotMessage, 0)
	err := GetDB().Table("bot_messages").Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return r
}

func GetAllMessageCommands() *[]string {
	ms := GetBotMessages()
	if ms != nil {
		cs := make([]string, 0)
		for _, v := range ms {
			cs = append(cs, v.Command)
		}
		return &cs
	}
	return nil
}

func GetBotMessageById(id string) *BotMessage {

	r := make([]*BotMessage, 0)
	err := GetDB().Table("bot_messages").Where("message_id = ?", id).Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return r[0]
}

func GetCommandsByMessage(command string) *BotMessage {

	r := make([]*BotMessage, 0)
	err := GetDB().Table("bot_messages").Where("command = ?", command).Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return r[0]
}
