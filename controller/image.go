package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type image struct {
	Created int64 `json:"created,omitempty"`
	Data    []struct {
		Url string `json:"url,omitempty"`
	} `json:"data,omitempty"`
}

func Image(context *gin.Context) {
	question := context.PostForm("question")
	key := context.PostForm("key")

	client := &http.Client{}
	body := make(map[string]interface{})

	body["prompt"] = question
	body["n"] = 1
	body["size"] = "256x256"

	b, e := json.Marshal(body)

	if e != nil {
		fmt.Println("err")
	}

	bReader := bytes.NewReader(b)

	res, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bReader)

	if err != nil {
		fmt.Println("error")
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", "Bearer "+key)
	res.Header = header

	rep, _ := client.Do(res)
	defer rep.Body.Close()

	if rep.StatusCode == http.StatusOK {
		decode := json.NewDecoder(rep.Body)

		r := &image{}
		err := decode.Decode(r)

		if err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}

		context.JSON(http.StatusOK, r.Data)

		return
	}

	context.JSON(http.StatusBadRequest, rep.Body)
}
