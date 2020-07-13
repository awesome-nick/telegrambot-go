package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	ImageCaption string `json:"caption"`
	Filepath     string `json:"filepath"`
	Command      string `json:"command" gorm:"unique;not null"`
}

func (i *Image) Create() {
	var img Image
	if err := GetDB().Where("command = ?", i.Command).First(&img).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			GetDB().Create(&i)
		}
	} else {
		GetDB().Model(&img).Where("command = ?", i.Command).Update(&i)
	}
}

func (i *Image) Delete() bool {

	var img Image
	if err := GetDB().Where("command = ?", i.Command).First(&img).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
	}
	GetDB().Unscoped().Delete(&i)
	return true
}

func GetAllImages() []*Image {

	records := make([]*Image, 0)
	err := GetDB().Table("images").Find(&records).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return records
}

func GetImageCommands() map[string]*Image {
	imgs := GetAllImages()
	m := make(map[string]*Image)

	for _, i := range imgs {
		m[i.Command] = i
	}

	return m
}

func GetAllImageCommands() *[]string {
	is := GetImageCommands()
	cs := make([]string, 0)
	for c := range is {
		cs = append(cs, c)
	}
	return &cs
}

func GetImgByCommand(c string) (*Image, bool) {
	m := GetImageCommands()

	if i, ok := m[c]; ok {
		return i, true
	}
	return nil, false
}
