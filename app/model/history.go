package model

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
)

type History struct {
	BaseModel
	Id        uint32 `gorm:"column:id;type:INT UNSIGNED;AUTO_INCREMENT;NOT NULL"`
	Uid       uint32 `gorm:"column:uid;type:INT UNSIGNED;NOT NULL"`
	Content   string `gorm:"column:content;type:TEXT;NOT NULL"`
	IsDeleted int8   `gorm:"column:is_deleted;type:TINYINT(1);NOT NULL"`
}

func (*History) Create(history *History) uint32 {
	db := global.DataBase.Create(history)

	if db.Error != nil {
		fmt.Println(db.Error.Error())

		return 0
	}

	return history.Id
}

func (*History) TopTen(condition map[string]interface{}, res []*History) []*History {
	global.DataBase.Where(condition).Offset(0).Limit(8).Order("id desc").Find(&res)

	return res
}

func (*History) Update(condition map[string]interface{}, updates map[string]interface{}) bool {
	db := global.DataBase.Model(&History{}).Where(condition).Updates(updates)

	if db.Error != nil {
		fmt.Println(db.Error.Error())

		return false
	}

	return true
}
