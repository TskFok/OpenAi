package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Upload(ctx *gin.Context) {
	file := ctx.PostForm("file")
	fmt.Println(file)
}
