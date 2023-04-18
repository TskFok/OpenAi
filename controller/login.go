package controller

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/model"
	"github.com/TskFok/OpenAi/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")

	if !tool.IsEmail(email) {
		ctx.JSON(http.StatusNotFound, "邮箱格式错误")
		return
	}
	password := ctx.PostForm("password")

	em := make(map[string]string)
	em["email"] = email
	mu := &model.User{}
	u := mu.Find(em)

	if u == nil {
		ctx.JSON(http.StatusNotFound, "邮箱不存在或密码错误")
		return
	}

	if u.Status != 1 {
		ctx.JSON(http.StatusNotFound, "用户状态错误")
		return
	}

	nPassword := tool.Password(password, u.Salt)

	if nPassword != u.Password {
		ctx.JSON(http.StatusNotFound, "邮箱不存在或密码错误")
		return
	}

	token, err := tool.JwtToken(u.Id, u.Email)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "生成token错误")
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func Register(ctx *gin.Context) {
	email := ctx.PostForm("email")

	if !tool.IsEmail(email) {
		ctx.JSON(http.StatusNotFound, "邮箱格式错误")
		return
	}
	password := ctx.PostForm("password")

	passwordErr := tool.CheckPasswordLever(password)

	if nil != passwordErr {
		ctx.JSON(http.StatusNotFound, passwordErr.Error())
		return
	}

	fmt.Println(email)
	fmt.Println(password)

	em := make(map[string]string)
	em["email"] = email
	mu := &model.User{}
	u := mu.Find(em)

	if u != nil {
		ctx.JSON(http.StatusBadRequest, "注册失败,用户邮箱已存在")
		return
	}

	salt := tool.UUID()
	nPassword := tool.Password(password, salt)

	mu.Status = 0
	mu.Salt = salt
	mu.Password = nPassword
	mu.Email = email

	uid := mu.Create(mu)

	if uid != 0 {
		token, err := tool.JwtToken(uid, email)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, "生成token错误")
			return
		}

		ctx.JSON(http.StatusOK, token)
		return
	}

	ctx.JSON(http.StatusBadRequest, "注册失败")
	return
}
