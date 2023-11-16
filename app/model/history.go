package model

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
)

type history struct {
	BaseModel
	Id        uint32 `gorm:"column:id;type:INT UNSIGNED;AUTO_INCREMENT;NOT NULL"`
	Uid       uint32 `gorm:"column:uid;type:INT UNSIGNED;NOT NULL"`
	Content   string `gorm:"column:content;type:TEXT;NOT NULL"`
	IsDeleted int8   `gorm:"column:is_deleted;type:TINYINT(1);NOT NULL"`
}

func NewHistory() *history {
	return new(history)
}

func NewHistorySlice(len int) []*history {
	return make([]*history, len)
}

func (*history) Create(history *history) uint32 {
	db := global.DataBase.Create(history)

	if db.Error != nil {
		fmt.Println(db.Error.Error())

		return 0
	}

	return history.Id
}

func (*history) TopTen(condition map[string]interface{}, res []*history) []*history {
	global.DataBase.Where(condition).Offset(0).Limit(8).Order("id desc").Find(&res)

	return res
}

func (*history) Update(condition map[string]interface{}, updates map[string]interface{}) bool {
	db := global.DataBase.Model(&history{}).Where(condition).Updates(updates)

	if db.Error != nil {
		fmt.Println(db.Error.Error())

		return false
	}

	return true
}
