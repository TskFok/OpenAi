package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/app/model"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Stream(ctx *gin.Context) {
	w := ctx.Writer

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, ok := w.(http.Flusher)

	if !ok {
		log.Panic("server not support") //浏览器不兼容
	}
	que := ctx.PostForm("question")

	hm := &model.History{}
	userId, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, "用户不存在")
		return
	}

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

	var a []stream
	jsonErr := json.Unmarshal([]byte(que), &a)

	if jsonErr != nil {
		fmt.Println(jsonErr.Error())
	}
	length := len(a)

	hm.Uid = userId.(uint32)
	hm.Content = a[length-1 : length][0].Value
	hm.IsDeleted = 0
	hmId := hm.Create(hm)

	if hmId == 0 {
		ctx.JSON(http.StatusBadRequest, "记录历史记录错误")
		return
	}

	megList := make([]openai.ChatCompletionMessage, 0)
	for i, v := range a {
		if i == 0 {
			continue
		}

		role := openai.ChatMessageRoleUser
		if v.Type == "answer" {
			role = openai.ChatMessageRoleAssistant
		}
		megList = append(megList, openai.ChatCompletionMessage{
			Role:    role,
			Content: v.Value,
		})
	}

	if len(megList) > 7 {
		megList = megList[len(megList)-7:]
	}

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 3000,
		Messages:  megList,
		Stream:    true,
	}
	stream, err := c.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		return
	}
	defer stream.Close()

	msg := make(chan interface{})
	stop := make(chan interface{})

	go func() {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				stop <- "stop"
				break
			}

			if err != nil {
				ctx.Done()
				stop <- "error"
				break
			}

			msg <- response.Choices[0].Delta.Content
		}
	}()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case message, ok := <-msg:
			if ok {
				w.Write([]byte(message.(string)))
			}
			return ok
		case <-stop:
			return false
		}
	})
}
