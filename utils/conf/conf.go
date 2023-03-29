package conf

import (
	"bytes"
	_ "embed"
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

}
