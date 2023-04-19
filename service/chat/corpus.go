package chat

import (
	"encoding/json"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/TskFok/OpenAi/utils/curl"
	"net/http"
)

type res struct {
	Object string `json:"object,omitempty"`
	Data   []struct {
		Object    string    `json:"object,omitempty"`
		Index     int       `json:"index,omitempty"`
		Embedding []float64 `json:"embedding,omitempty"`
	} `json:"data,omitempty"`
	Model string `json:"model,omitempty"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens,omitempty"`
		TotalTokens  int `json:"total_tokens,omitempty"`
	} `json:"usage"`
	Corpus string
}

func File() {
	//使用语料库
	body := make(map[string]interface{})
	body["model"] = "text-embedding-ada-002"
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+global.OpenAiToken)

	keys := cache.Keys("embeding:*")

	for _, v := range keys {
		cache.Del(v)
	}

	for i, v := range global.Corpus {
		body["input"] = v

		req := &res{}
		httpStatus := curl.Post("https://api.openai.com/v1/embeddings", body, header, req)

		if httpStatus != http.StatusOK {
			fmt.Println("查询失败")
		}

		req.Corpus = v.(string)

		by, err := json.Marshal(req)

		if err != nil {
			fmt.Println(err.Error())
		}

		cache.Set(fmt.Sprintf("embeding:%v", i), string(by), 3600)
	}

}
