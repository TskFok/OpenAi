package me

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/TskFok/OpenAi/utils/curl"
	"math"
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

type rs struct {
	Data   []float64
	Corpus string
}

func file(question string) (string, error) {
	//使用语料库
	body := make(map[string]interface{})
	body["model"] = "text-embedding-ada-002"
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+global.OpenAiToken)

	body["input"] = question

	requestion := &res{}
	httpStatus := curl.Post("https://api.openai.com/v1/embeddings", body, header, requestion)

	if httpStatus != http.StatusOK {
		return "", errors.New("查询失败")
	}
	redisKeys := cache.Keys("embeding:*")

	var lastFa float64 = -10
	var corpusDetail string

	var fa3 float64 = 0
	for _, v := range requestion.Data[0].Embedding {
		fa3 += math.Pow(v, 2)
	}
	for _, v := range redisKeys {
		rr := &rs{}
		val := cache.Get(v)
		err := json.Unmarshal([]byte(val), rr)

		if err != nil {
			fmt.Println(err.Error())
		}
		var fa float64 = 0
		var fa2 float64 = 0

		for i, v2 := range rr.Data {
			fa2 += math.Pow(v2, 2)
			fa += v2 * requestion.Data[0].Embedding[i]
		}

		fi := fa / (math.Sqrt(fa2) * math.Sqrt(fa3))

		if fi > lastFa {
			lastFa = fi
			corpusDetail = rr.Corpus
		}
	}
	return "Answer the question as truthfully as possible using the provided context, and if the answer is not contained within the text below, say \"I don't know.\"\\n\\nContext:\\n" + corpusDetail + "\n\n Q: ", nil
}

//"We have provided context information below: \n"
//"---------------------\n"
//"{context_str}\n"
//"---------------------\n"
//"Given this information, Please answer my question in the same language that I used to ask you.\n"
//"Please answer the question: {query_str}\n"
