package controller

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/model"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func DeleteHistory(ctx *gin.Context) {
	uid, exists := ctx.Get("user_id")

	if !exists {
		fmt.Println("用户信息不存在")

		return
	}
	uidd := uid.(uint32)

	hm := &model.History{}

	condition := make(map[string]interface{})
	condition["uid"] = uidd

	updates := make(map[string]interface{})
	updates["is_deleted"] = 1
	hm.Update(condition, updates)

	iUid := int(uidd)
	sUid := strconv.Itoa(iUid)
	cache.Del("history:" + sUid)

	ctx.JSON(http.StatusOK, "success")
}

func HistoryList(ctx *gin.Context) {
	uid, exists := ctx.Get("user_id")

	if !exists {
		fmt.Println("用户信息不存在")
	}

	condition := make(map[string]interface{})
	condition["uid"] = uid.(uint32)
	condition["is_deleted"] = 0

	historyMap := make([]*model.History, 8)

	hm := &model.History{}
	hm.TopTen(condition, historyMap)

	ctx.JSON(http.StatusOK, historyMap)
}
