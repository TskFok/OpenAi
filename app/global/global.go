package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var TlsCert string
var TlsKey string
var AppMode string
var RedisHost string
var RedisUser string
var RedisPassword string
var RedisClient *redis.Client
var DataBase *gorm.DB
var WechatAppid string
var WechatSecret string
var MysqlDsn string
var MysqlPrefix string
