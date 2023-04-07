package conf

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/spf13/viper"
)

//go:embed conf.yaml
var conf []byte

func InitConfig() {
	viper.SetConfigType("yaml")

	err := viper.ReadConfig(bytes.NewReader(conf))

	if nil != err {
		panic(err)
	}

	global.TlsCert = viper.Get("tls.cert").(string)
	global.TlsKey = viper.Get("tls.key").(string)
	global.AppMode = viper.Get("app.mode").(string)
	global.RedisPassword = viper.Get("redis.password").(string)
	global.RedisHost = viper.Get("redis.host").(string)

	fmt.Println("获取配置:" + viper.Get("wechat.appid").(string))
	global.WechatAppid = viper.Get("wechat.appid").(string)
	global.WechatSecret = viper.Get("wechat.secret").(string)

}
