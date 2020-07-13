package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type APIUser struct {
	ID       uint64 `json:"id" gorm:"unique;not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password"`
}

func (au *APIUser) Create() bool {

	var n APIUser
	if err := GetDB().Where("username = ?", au.Username).First(&n).Error; err == nil {
		return false
	}
	au.Password = hashPasswd(au.Password)
	GetDB().Create(au)
	return true
}

func (au *APIUser) Update() bool {

	var u APIUser
	if err := GetDB().Where("username = ?", au.Username).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
	}
	GetDB().Model(&u).Where("username = ?", au.Username).Update(&au)
	return true
}

func (au *APIUser) Delete() bool {

	var u APIUser
	if err := GetDB().Where("username = ?", au.Username).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
	}
	GetDB().Model(&u).Delete(&au)
	return true
}

func CreateAdmin(password string) {
	p := hashPasswd(password)
	if p == "" {
		log.Println("Smth wrong with hashing password!")
		return
	}

	var admin = APIUser{
		Username: "admin",
		Password: p,
	}
	_ = admin.Create()

}

func GetAPIUsers() []*APIUser {

	r := make([]*APIUser, 0)
	err := GetDB().Table("api_users").Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return r
}

func GetAPIUserByUserName(username string) *APIUser {

	r := make([]*APIUser, 0)
	err := GetDB().Table("api_users").Where("username = ?", username).Find(&r).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(r) > 0 {
		return r[0]
	}
	return nil
}

func hashPasswd(p string) string {

	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(h)
}

func CompareHashAndPasswd(h, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	if err == nil {
		return true
	} else {
		return false
	}

}
