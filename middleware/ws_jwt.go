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

type wsu struct {
	Id     uint32
	Email  string
	Status int8
}

func Validate(ctx *gin.Context) int {
	token := ctx.GetHeader("Sec-WebSocket-Protocol")
	claims, tokenErr := tool.TokenInfo(token)

	if nil != tokenErr {
		return http.StatusUnauthorized
	}

	builder := strings.Builder{}
	builder.WriteString("user:info:")
	builder.WriteString(strconv.FormatUint(uint64(claims.Uid), 10))
	key := builder.String()

	if cache.Has(key) {
		user := &wsu{}
		jsonErr := json.Unmarshal([]byte(cache.Get(key)), &user)

		if nil != jsonErr {
			return http.StatusUnauthorized
		}

		ctx.Set("user_status", user.Status)
		ctx.Set("email", user.Email)
		ctx.Set("user_id", user.Id)
	} else {
		userModel := model.NewUser()
		condition := make(map[string]interface{})
		condition["id"] = claims.Uid
		user := userModel.Find(condition)

		if user == nil {
			return http.StatusUnauthorized
		}

		if user.Status != 1 {
			return http.StatusUnauthorized
		}

		ctx.Set("user_status", user.Status)
		ctx.Set("email", user.Email)
		ctx.Set("user_id", user.Id)

		res, jsonErr := json.Marshal(user)

		if nil != jsonErr {
			return http.StatusUnauthorized
		}
		cache.Set(key, string(res), 3600)
	}

	return http.StatusOK
}
