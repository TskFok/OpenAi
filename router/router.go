package router

import (
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/app/websockets/chat"
	"github.com/TskFok/OpenAi/app/websockets/me"
	"github.com/TskFok/OpenAi/controller"
	"github.com/TskFok/OpenAi/middleware"
	"github.com/TskFok/OpenAi/public/html"
	"github.com/TskFok/OpenAi/service/applet"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

var Handle *gin.Engine

func InitRouter() {
	go chat.WebsocketManager.Start()

	gin.SetMode(global.AppMode)

	Handle = gin.New()
	Handle.Use(gin.Recovery())
	Handle.Use(gin.Logger())
	Handle.Use(middleware.Cors())

	Handle.POST("/login", controller.Login)
	Handle.POST("/register", controller.Register)

	wsGroup := Handle.Group("/gpt")
	{
		wsGroup.GET("/:channel", chat.WebsocketManager.WsClient)
	}

	wsGroupMe := Handle.Group("/me")
	{
		wsGroupMe.GET("/:channel", me.WebsocketManager.WsClient)
	}

	Handle.Use(middleware.Jwt())
	Handle.GET("/history", controller.HistoryList)
	Handle.DELETE("/history", controller.DeleteHistory)
	Handle.POST("/chat", controller.Chat)
	Handle.POST("/chat2", controller.Chat2)
	Handle.POST("/image", controller.Image)
	Handle.POST("/upload", controller.Upload)
	Handle.GET("/stream", controller.ChatStream)

	Handle.SetHTMLTemplate(html.GetQuestionTemplate())

	Handle.GET("/scan", func(c *gin.Context) {
		c.HTML(http.StatusOK, "scan.html", gin.H{})
	})

	Handle.GET("/chat-web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "question.html", gin.H{})
	})

	Handle.GET("/image", func(c *gin.Context) {
		c.HTML(http.StatusOK, "image.html", gin.H{})
	})

	Handle.GET("/wx-applet", func(c *gin.Context) {
		c.HTML(http.StatusOK, "wxApplet.html", gin.H{
			"code": template.URL(applet.GetCode()),
		})
	})

}
