package bootstrap

import (
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/utils/cache"
	"github.com/TskFok/OpenAi/utils/conf"
)

func Init() {
	conf.InitConfig()

	global.RedisClient = cache.InitRedis()
}
