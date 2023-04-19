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

type stream struct {
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

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
