package vicuna

import (
	"bytes"
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
	Data   []float64 `json:"data,omitempty"`
	Corpus string    `json:"corpus,omitempty"`
	Pid    int       `json:"pid,omitempty"`
	Id     int       `json:"id,omitempty"`
}

type rMap struct {
	Corpus string
	Fi     float64
}

func file(question, id string) (string, error) {
	//使用语料库
	body := make(map[string]interface{})
	body["model"] = "vicuna-13b"
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "")

	body["input"] = question

	requestion := &res{}
	httpStatus := curl.Post(global.VicunaUrl+"/embeddings", body, header, requestion)

	if httpStatus != http.StatusOK {
		return "", errors.New("查询失败")
	}
	redisKeys := cache.Keys("embeding_new:*")

	var lastFa float64 = -10

	var fa3 float64 = 0
	for _, v := range requestion.Data[0].Embedding {
		fa3 += math.Pow(v, 2)
	}

	ppid := 0
	iid := 0

	fiMap := make(map[string]rMap)
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

		//fmt.Printf("fa:%f fa2:%f fa3:%f fi:%f", fa, fa2, fa3, fi)
		//fmt.Println(rr.Corpus)

		rss := rMap{}
		rss.Fi = fi
		rss.Corpus = rr.Corpus
		fiMap[v] = rss

		if fi > lastFa {
			lastFa = fi
			ppid = rr.Pid
			iid = rr.Id
		}
	}

	//提问集体信息，语料为基础类目的所有语料+该类目下相似程度较高的语料
	key := fmt.Sprintf("embeding_new:0:%d", ppid)
	if ppid == 0 {
		//提问基础信息，语料为基础类目下的所有语料
		key = fmt.Sprintf("embeding_new:0:%d", iid)
	}

	baseRr := &rs{}
	val := cache.Get(key)
	err := json.Unmarshal([]byte(val), baseRr)

	if err != nil {
		fmt.Println(err.Error())
	}

	useKeys := cache.Keys(fmt.Sprintf("embeding_new:%d:*", ppid))
	if ppid == 0 {
		useKeys = cache.Keys(fmt.Sprintf("embeding_new:%d:*", iid))
	}

	for _, v := range useKeys {
		fMaps := fiMap[v]

		if fMaps.Fi > 0.5 {
			cache.ZAdd(fmt.Sprintf("embeding_new_list:%s", id), fMaps.Fi, fMaps.Corpus)
		}
	}

	bf := buffer.Buffer{}
	bf.WriteString("We have provided context information below: \n")
	bf.WriteString("---------------------\n")
	bf.WriteString(baseRr.Corpus)
	bf.WriteString("\n")

	//取最接近的三十个语料合并作为一个提示,语料越细越好
	for _, v := range cache.ZRangeAndRemove("embeding_new_list:"+id, "0", "1", 0, 30) {
		bf.WriteString(v)

		if bytes.Count(bf.Bytes(), nil) > 600 {
			break
		}
		bf.WriteString(",")
	}

	if len(useKeys) != 0 {
		bf.WriteString("\n")
	}

	bf.WriteString("---------------------\n")
	bf.WriteString("Given this information, Please answer my question in the same language that I used to ask you.\n")
	bf.WriteString("and if the answer is not contained within the text below, say \"我不知道.\"")
	bf.WriteString("Please answer the question: \n")

	return bf.String(), nil
}

//"We have provided context information below: \n"
//"---------------------\n"
//"{context_str}\n"
//"---------------------\n"
//"Given this information, Please answer my question in the same language that I used to ask you.\n"
//"Please answer the question: {query_str}\n"
