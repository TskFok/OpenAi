package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"net/url"
)

func ChatStream(ctx *gin.Context) {
	w := ctx.Writer

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, ok := w.(http.Flusher)

	if !ok {
		log.Panic("server not support") //浏览器不兼容
	}
	que := ctx.Query("question")
	config := openai.DefaultConfig(global.OpenAiToken)
	//使用warp代理,不使用代理 cai := openai.NewClient(send.Key)
	proxyUrl, err := url.Parse("http://127.0.0.1:40000")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	c := openai.NewClientWithConfig(config)

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 3000,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: que,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		return
	}
	defer stream.Close()

	msg := make(chan interface{})
	stop := make(chan interface{})

	//开个协程不断获取stream数据
	go func(ctx *gin.Context) {
		for {
			response, err := stream.Recv()

			//为空发送stop信号
			if errors.Is(err, io.EOF) {
				stop <- "stop"
				break
			}

			//报错发送error信号
			if err != nil {
				fmt.Println(err.Error())
				stop <- "error"
				break
			}

			content := response.Choices[0].Delta.Content

			//英文回答开头第一个空白字符串特殊处理 返回一个不是空白的字符串让前端进行替换
			if len(content) > 1 && content[0:1] == " " {
				msg <- "<<emptystring>>"
			}

			msg <- response.Choices[0].Delta.Content
		}
	}(ctx)

	//流式返回数据
	ctx.Stream(func(w io.Writer) bool {
		select {
		case message, ok := <-msg:
			if ok {
				ctx.SSEvent("message", message)
			}
			return ok
		case tp := <-stop:
			ctx.SSEvent(tp.(string), "stop")
			return false
		case <-ctx.Done():
			ctx.SSEvent("error", "done")
			return false
		}
	})

}
