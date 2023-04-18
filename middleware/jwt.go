package middleware

import (
	"encoding/json"
	"github.com/TskFok/OpenAi/app/model"
	"github.com/TskFok/OpenAi/tool"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type u struct {
	Id     uint32
	Email  string
	Status int8
}

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		sec := ctx.GetHeader("Sec-WebSocket-Protocol")

		if sec != "" {
			token = sec
		}

		claims, tokenErr := tool.TokenInfo(token)

		if nil != tokenErr {
			ctx.JSON(http.StatusUnauthorized, tokenErr.Error())
			ctx.Abort()
			return
		}

		builder := strings.Builder{}
		builder.WriteString("user:info:")
		builder.WriteString(strconv.FormatUint(uint64(claims.Uid), 10))
		key := builder.String()

		if cache.Has(key) {
			user := &u{}
			jsonErr := json.Unmarshal([]byte(cache.Get(key)), &user)

			if nil != jsonErr {
				ctx.JSON(http.StatusUnauthorized, jsonErr.Error())
				ctx.Abort()
				return
			}

			ctx.Set("user_status", user.Status)
			ctx.Set("email", user.Email)
			ctx.Set("user_id", user.Id)
		} else {
			userModel := &model.User{}
			condition := make(map[string]interface{})
			condition["id"] = claims.Uid
			user := userModel.Find(condition)

			if user == nil {
				ctx.JSON(http.StatusUnauthorized, "用户不存在")
				ctx.Abort()
				return
			}

			if user.Status != 1 {
				ctx.JSON(http.StatusUnauthorized, "用户状态错误")
				ctx.Abort()
				return
			}

			ctx.Set("user_status", user.Status)
			ctx.Set("email", user.Email)
			ctx.Set("user_id", user.Id)

			res, jsonErr := json.Marshal(user)

			if nil != jsonErr {
				ctx.JSON(http.StatusUnauthorized, jsonErr.Error())
				ctx.Abort()
				return
			}
			cache.Set(key, string(res), 3600)
		}

		ctx.Next()
	}
}
