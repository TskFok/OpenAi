/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/TskFok/OpenAi/utils/curl"
	"github.com/spf13/cobra"
)

type tokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Errcode     string `json:"errcode,omitempty"`
}

// wxTokenCmd represents the test command
var wxTokenCmd = &cobra.Command{
	Use:   "wx:token",
	Short: "微信token",
	Long:  `获取微信token,并存入cache`,
	Run: func(cmd *cobra.Command, args []string) {
		tkq := &tokenResponse{}
		curl.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+global.WechatAppid+"&secret="+global.WechatSecret, tkq)

		if tkq.Errcode == "" {
			cache.Set("wx:token", tkq.AccessToken, 3600)
		} else {
			fmt.Println(tkq)
		}
	},
}

func init() {
	rootCmd.AddCommand(wxTokenCmd)
}
