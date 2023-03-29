package router

import (
	"github.com/TskFok/OpenAi/app/websockets"
	"github.com/TskFok/OpenAi/controller"
	"github.com/TskFok/OpenAi/middleware"
	"github.com/TskFok/OpenAi/public/html"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Handle *gin.Engine

func InitRouter() {
	go websockets.WebsocketManager.Start()
	go websockets.WebsocketManager.SendService()

	gin.SetMode(gin.ReleaseMode)

	Handle = gin.New()
	Handle.Use(gin.Recovery())
	Handle.Use(gin.Logger())
	Handle.Use(middleware.Cors())

	Handle.POST("/chat", controller.Chat)
	Handle.POST("/chat2", controller.Chat2)
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

	wsGroup := Handle.Group("/gpt")
	{
		wsGroup.GET("/:channel", websockets.WebsocketManager.WsClient)
	}
}
