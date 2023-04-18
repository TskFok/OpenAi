package model

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
)

type User struct {
	BaseModel
	Id       uint32 `gorm:"column:id;type:INT UNSIGNED;AUTO_INCREMENT;NOT NULL"`
	Email    string `gorm:"column:email;type:VARCHAR(255);NOT NULL"`
	Salt     string `gorm:"column:salt;type:VARCHAR(50);NOT NULL"`
	Password string `gorm:"column:password;type:VARCHAR(255);NOT NULL"`
	Status   int8   `gorm:"column:status;type:TINYINT(1);NOT NULL"`
}

func (*User) Find(where interface{}) (u *User) {
	db := global.DataBase.Where(where).First(&u)

	if db.Error != nil {
		return nil
	}

	return u
}

func (*User) Create(user *User) uint32 {
	db := global.DataBase.Create(user)

	if db.Error != nil {
		fmt.Println(db.Error.Error())
		return 0
	}

	return user.Id
}
