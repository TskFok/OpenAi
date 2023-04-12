package applet

import (
	"encoding/base64"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/TskFok/OpenAi/utils/curl"
	"net/http"
)

type postBody struct {
	Scene      string `json:"scene,omitempty"`
	Page       string `json:"page,omitempty"`
	Width      int    `json:"width,omitempty"`
	CheckPath  bool   `json:"check_path"`
	EnvVersion string `json:"env_version,omitempty"`
}

func GetCode() string {
	body := &postBody{
		Scene:      "1",
		Page:       "pages/upload/upload",
		Width:      500,
		CheckPath:  false,
		EnvVersion: "trial",
	}

	header := http.Header{}

	image := curl.PostAll("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+cache.Get("wx:token"),
		body, header)

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
}
