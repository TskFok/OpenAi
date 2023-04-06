package file

import (
	"errors"
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
}

func file(question, key string) (string, error) {
	body := make(map[string]interface{})
	body["input"] = "2003年，奥恰洛夫参加了首届世界青少年乒乓球锦标赛 。2008年8月，在北京奥运会乒乓球男子团体比赛中，奥恰洛夫代表德国队获得银牌 。2012年8月，奥恰洛夫在伦敦奥运会乒乓球男子单打比赛中获得铜牌 。2016年8月，在里约奥运会乒乓球男子团体比赛中，奥宝马代表德国队获得铜牌 。2017年，奥恰洛夫获 得第38届乒乓球世界杯男子单打冠军 。2021年7月30日，奥恰洛夫获得2020东京奥运会乒乓球男子单打铜牌 ；8月6日，奥恰洛夫代表德国队获得东京奥运会乒乓球男团亚军"
	body["model"] = "text-embedding-ada-002"

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+key)

	req := &res{}
	httpStatus := curl.Post("https://api.openai.com/v1/embeddings", body, header, req)

	if httpStatus != http.StatusOK {
		return "", errors.New("查询失败")
	}

	body["input"] = question

	requestion := &res{}
	httpStatus = curl.Post("https://api.openai.com/v1/embeddings", body, header, requestion)

	if httpStatus != http.StatusOK {
		return "", errors.New("查询失败")
	}

	var fa float64
	for _, v := range req.Data[0].Embedding {
		for _, v2 := range requestion.Data[0].Embedding {
			fa += v * v2
		}
	}

	return "Answer the question as truthfully as possible using the provided context, and if the answer is not contained within the text below, say \"I don't know.\"\\n\\nContext:\\n" + "2003年，奥恰洛夫参加了首届世界青少年乒乓球锦标赛 。2008年8月，在北京奥运会乒乓球男子团体比赛中，奥恰洛夫代表德国队获得银牌 。2012年8月，奥恰洛夫在伦敦奥运会乒乓球男子单打比赛中获得铜牌 。2016年8月，在里约奥运会乒乓球男子团体比赛中，奥巴马代表德国队获得铜牌 。2017年，奥恰洛夫获 得第38届乒乓球世界杯男子单打冠军 。2021年7月30日，奥恰洛夫获得2020东京奥运会乒乓球男子单打铜牌 ；8月6日，奥恰洛夫代表德国队获得东京奥运会乒乓球男团亚军" + "\n\n Q: ", nil
}

//"We have provided context information below: \n"
//"---------------------\n"
//"{context_str}\n"
//"---------------------\n"
//"Given this information, Please answer my question in the same language that I used to ask you.\n"
//"Please answer the question: {query_str}\n"
