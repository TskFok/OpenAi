package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/app/model"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap/buffer"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type stream struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Setup struct {
	Token           int     `json:"token,omitempty"`
	Temperature     float32 `json:"temperature,omitempty"`
	PresencePenalty float32 `json:"presence_penalty,omitempty"`
	HistoryNum      int     `json:"history_num,omitempty"`
}

func ChatSse(ctx *gin.Context) {
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
	setup := ctx.Query("setup")

	hm := &model.History{}
	userId, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, "用户不存在")

		return
	}

	sp := &Setup{}
	err := json.Unmarshal([]byte(setup), sp)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "配置错误")

		return
	}

	userIid := userId.(uint32)

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

	hm.Uid = userIid
	hm.Content = que
	hm.IsDeleted = 0
	hmId := hm.Create(hm)

	if hmId == 0 {
		ctx.JSON(http.StatusBadRequest, "记录历史记录错误")
		return
	}

	iUid := int(userIid)
	sUid := strconv.Itoa(iUid)
	history := cache.Get("history:" + sUid)
	streams := make([]stream, 10)

	if history != "" {
		err = json.Unmarshal([]byte(history), &streams)

		if err != nil {
			fmt.Println(err)
		}
	}

	megList := make([]openai.ChatCompletionMessage, 0)
	for _, v := range streams {
		if v.Role == "" {
			break
		}
		role := openai.ChatMessageRoleUser
		if v.Role == "answer" {
			role = openai.ChatMessageRoleAssistant
		}
		megList = append(megList, openai.ChatCompletionMessage{
			Role:    role,
			Content: v.Content,
		})
	}

	megList = append(megList, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: que,
	})

	if len(megList) > sp.HistoryNum {
		megList = megList[len(megList)-sp.HistoryNum:]
	}

	req := openai.ChatCompletionRequest{
		Model:           openai.GPT3Dot5Turbo,
		MaxTokens:       sp.Token,
		Messages:        megList,
		Stream:          true,
		Temperature:     sp.Temperature,
		PresencePenalty: sp.PresencePenalty,
	}

	stream, err := c.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stream.Close()

	msg := make(chan interface{})
	stop := make(chan interface{})

	//开个协程不断获取stream数据
	go func() {
		bf := buffer.Buffer{}
		for {
			response, err := stream.Recv()

			//为空发送stop信号
			if errors.Is(err, io.EOF) {
				writeAnswer(bf.String(), "history:"+sUid, megList)
				stop <- "stop"
				break
			}

			//报错发送error信号
			if err != nil {
				writeAnswer(bf.String(), "history:"+sUid, megList)
				fmt.Println(err.Error())
				stop <- "error"
				break
			}

			content := response.Choices[0].Delta.Content

			//英文回答开头第一个空白字符串特殊处理 返回一个不是空白的字符串让前端进行替换
			if len(content) > 1 && content[0:1] == " " {
				msg <- "<<emptystring>>"
			}

			bf.WriteString(response.Choices[0].Delta.Content)
			msg <- response.Choices[0].Delta.Content
		}
	}()

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
		}
	})

}

func writeAnswer(answer, key string, megList []openai.ChatCompletionMessage) {
	megList = append(megList, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	})

	bt, err := json.Marshal(megList)

	if err != nil {
		fmt.Println(err.Error())
	}
	cache.Set(key, string(bt), 3600)
}

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
	setup := ctx.PostForm("setup")

	sp := &Setup{}
	err := json.Unmarshal([]byte(setup), sp)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "配置错误")

		return
	}

	hm := &model.History{}
	userId, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, "用户不存在")
		return
	}

	userIid := userId.(uint32)

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

	hm.Uid = userIid
	hm.Content = que
	hm.IsDeleted = 0
	hmId := hm.Create(hm)

	if hmId == 0 {
		ctx.JSON(http.StatusBadRequest, "记录历史记录错误")
		return
	}

	iUid := int(userIid)
	sUid := strconv.Itoa(iUid)
	history := cache.Get("history:" + sUid)
	streams := make([]stream, 10)

	if history != "" {
		err = json.Unmarshal([]byte(history), &streams)

		if err != nil {
			fmt.Println(err)
		}
	}

	megList := make([]openai.ChatCompletionMessage, 0)
	for _, v := range streams {
		if v.Role == "" {
			break
		}

		role := openai.ChatMessageRoleUser
		if v.Role == "answer" {
			role = openai.ChatMessageRoleAssistant
		}
		megList = append(megList, openai.ChatCompletionMessage{
			Role:    role,
			Content: v.Content,
		})
	}

	megList = append(megList, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: que,
	})

	if len(megList) > sp.HistoryNum {
		megList = megList[len(megList)-sp.HistoryNum:]
	}

	req := openai.ChatCompletionRequest{
		Model:           openai.GPT3Dot5Turbo,
		MaxTokens:       sp.Token,
		Messages:        megList,
		Stream:          true,
		PresencePenalty: sp.PresencePenalty,
		Temperature:     sp.Temperature,
	}
	stream, err := c.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stream.Close()

	msg := make(chan interface{})
	stop := make(chan interface{})

	go func() {
		bf := buffer.Buffer{}
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				writeAnswer(bf.String(), "history:"+sUid, megList)
				stop <- "stop"
				break
			}

			if err != nil {
				writeAnswer(bf.String(), "history:"+sUid, megList)
				ctx.Done()
				stop <- "error"
				break
			}
			bf.WriteString(response.Choices[0].Delta.Content)
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
