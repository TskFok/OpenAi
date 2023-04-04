package router

import (
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/app/websockets/chat"
	"github.com/TskFok/OpenAi/app/websockets/file"
	"github.com/TskFok/OpenAi/controller"
	"github.com/TskFok/OpenAi/middleware"
	"github.com/TskFok/OpenAi/public/html"
	"github.com/gin-gonic/gin"
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

	Handle.POST("/chat", controller.Chat)
	Handle.POST("/chat2", controller.Chat2)
	Handle.POST("/image", controller.Image)
	Handle.GET("/stream", controller.ChatStream)

	Handle.SetHTMLTemplate(html.GetQuestionTemplate())

	Handle.GET("/chat-web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "question.html", gin.H{})
	})

	Handle.GET("/chat-web-sse", func(c *gin.Context) {
		c.HTML(http.StatusOK, "questionSse.html", gin.H{})
	})

	Handle.GET("/chat-web-ws", func(c *gin.Context) {
		c.HTML(http.StatusOK, "questionWs.html", gin.H{})
	})

	Handle.GET("/image", func(c *gin.Context) {
		c.HTML(http.StatusOK, "image.html", gin.H{})
	})

	Handle.GET("/chat-file", func(c *gin.Context) {
		c.HTML(http.StatusOK, "questionFileWs.html", gin.H{})
	})

	wsGroup := Handle.Group("/gpt")
	{
		wsGroup.GET("/:channel", chat.WebsocketManager.WsClient)
	}

	wsGroupFile := Handle.Group("/gpt-file")
	{
		wsGroupFile.GET("/:channel", file.WebsocketManager.WsClient)
	}
}
