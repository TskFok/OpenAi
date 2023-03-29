package controller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
)

func Chat2(ctx *gin.Context) {
	que := ctx.PostForm("question")
	//your key
	key := ctx.PostForm("key")
	c := openai.NewClient(key)

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

	go ctx.Stream(func(w io.Writer) bool {
		select {
		case message, ok := <-msg:
			if ok {
				w.Write([]byte(message.(string)))
			}
			return ok
		case <-ctx.Done():
			return false
		}
	})

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			ctx.Done()
			return
		}

		if err != nil {
			ctx.Done()
			return
		}

		msg <- response.Choices[0].Delta.Content
	}
}
