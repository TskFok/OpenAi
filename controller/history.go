package controller

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteHistory(ctx *gin.Context) {
	uid, exists := ctx.Get("user_id")

	if !exists {
		fmt.Println("用户信息不存在")
	}

	hm := &model.History{}

	condition := make(map[string]interface{})
	condition["uid"] = uid.(uint32)

	updates := make(map[string]interface{})
	updates["is_deleted"] = 1
	hm.Update(condition, updates)

	ctx.JSON(http.StatusOK, "success")
}
