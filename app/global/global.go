package global

import "github.com/redis/go-redis/v9"

var TlsCert string
var TlsKey string
var AppMode string
var RedisHost string
var RedisPassword string
var RedisClient *redis.Client
var WechatAppid string
var WechatSecret string
