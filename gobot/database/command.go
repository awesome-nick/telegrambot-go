package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Command struct {
	CommandName string `json:"command_name"`
	MessageId   string `json:"message_id" gorm:"unique;not null"`
}

func (command *Command) Create() {
	var cmd_holder Command
	if err := GetDB().Where("message_id = ?", command.MessageId).First(&cmd_holder).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			GetDB().Create(&command)
		}
	} else {
		GetDB().Model(&cmd_holder).Where("message_id = ?", command.MessageId).Update("command_name", command.CommandName)
	}
}

func GetAllCommands() []*Command {

	records := make([]*Command, 0)
	err := GetDB().Table("commands").Find(&records).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return records
}
