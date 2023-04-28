package me

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/TskFok/OpenAi/utils/curl"
	"go.uber.org/zap/buffer"
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

	bf := buffer.Buffer{}
	bf.WriteString("We have provided context information below: \n")
	bf.WriteString("---------------------\n")
	bf.WriteString(corpusDetail)
	bf.WriteString("\n")
	bf.WriteString("---------------------\n")
	bf.WriteString("Given this information, Please answer my question in the same language that I used to ask you.\n")
	bf.WriteString("and if the answer is not contained within the text below, say \"我不知道.\"")
	bf.WriteString("Please answer the question: ")

	return bf.String(), nil
}

//"We have provided context information below: \n"
//"---------------------\n"
//"{context_str}\n"
//"---------------------\n"
//"Given this information, Please answer my question in the same language that I used to ask you.\n"
//"Please answer the question: {query_str}\n"
