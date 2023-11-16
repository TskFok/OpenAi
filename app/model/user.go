package model

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
)

type user struct {
	BaseModel
	Id       uint32 `gorm:"column:id;type:INT UNSIGNED;AUTO_INCREMENT;NOT NULL"`
	Email    string `gorm:"column:email;type:VARCHAR(255);NOT NULL"`
	Salt     string `gorm:"column:salt;type:VARCHAR(50);NOT NULL"`
	Password string `gorm:"column:password;type:VARCHAR(255);NOT NULL"`
	Status   int8   `gorm:"column:status;type:TINYINT(1);NOT NULL"`
}

func NewUser() *user {
	return new(user)
}

func (*user) Find(where interface{}) (u *user) {
	db := global.DataBase.Where(where).First(&u)

	if db.Error != nil {
		return nil
	}

	return u
}

func (*user) Create(user *user) uint32 {
	db := global.DataBase.Create(user)

	if db.Error != nil {
		fmt.Println(db.Error.Error())
		return 0
	}

	return user.Id
}
