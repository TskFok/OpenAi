package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Header("Access-Control-Allow-Headers", "*")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}

		context.Next()
	}
}
