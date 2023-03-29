package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type gpt struct {
	Id      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Usage   usage  `json:"usage"`
	Choices []msg  `json:"choices,omitempty"`
}

type msg struct {
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason,omitempty"`
	Index        int64   `json:"index,omitempty"`
}

type message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type usage struct {
	PromptTokens     int64 `json:"prompt_tokens,omitempty"`
	CompletionTokens int64 `json:"completion_tokens,omitempty"`
	TotalTokens      int64 `json:"total_tokens,omitempty"`
}

type content struct {
	Text text
}

type text struct {
	Content string `json:"content,omitempty"`
}

type response struct {
	Msgtype string       `json:"msgtype,omitempty"`
	Text    responseText `json:"text"`
}

type responseText struct {
	Content string `json:"content,omitempty"`
}

func Chat(context *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(buf)

	ct := &content{}
	json.Unmarshal(buf[0:n], ct)

	client := &http.Client{}
	body := make(map[string]interface{})

	body["model"] = "gpt-3.5-turbo"
	msg := make([]map[string]string, 1)
	role := make(map[string]string)
	role["role"] = "user"
	role["content"] = ct.Text.Content
	msg[0] = role
	body["messages"] = msg
	body["temperature"] = 0.7
	body["max_tokens"] = 3000
	body["top_p"] = 1
	body["frequency_penalty"] = 0
	body["presence_penalty"] = 0.6
	body["stream"] = false

	b, e := json.Marshal(body)

	if e != nil {
		fmt.Println("err")
	}

	bReader := bytes.NewReader(b)

	res, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bReader)

	if err != nil {
		fmt.Println("error")
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer Your Key")
	res.Header = header

	rep, _ := client.Do(res)
	defer rep.Body.Close()

	if rep.StatusCode == http.StatusOK {
		decode := json.NewDecoder(rep.Body)

		r := &gpt{}
		decode.Decode(r)

		fmt.Println(r)
		respon := &response{
			Msgtype: "text",
			Text: responseText{
				Content: r.Choices[0].Message.Content,
			},
		}

		context.JSON(http.StatusOK, respon)
	} else {
		fmt.Println(rep.Status)
		fmt.Println(rep.Body)
		respon := &response{
			Msgtype: "text",
			Text: responseText{
				Content: rep.Status,
			},
		}

		context.JSON(http.StatusOK, respon)
	}
}
