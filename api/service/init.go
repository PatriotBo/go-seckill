package service

import (
	"github.com/astaxie/beego/config"
	log "github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

var (
	Config      config.Configer
	RedisClient *redis.Client
)

func Init() {
	initConfig()
	initRedis()
}

// 初始化redis
func initRedis() {
	addr := Config.String("redis.addr")
	password := Config.String("redis.password")
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := c.Ping().Result()
	if err != nil {
		log.Error("Connect redis failed. Error : %v", err)
		panic(err.Error())
	}
	RedisClient = c
}

//初始化配置
func initConfig() {
	iniConf, err := config.NewConfig("ini", "/config/api.conf")
	if err != nil {
		log.Error("init config failed:", err)
		panic("init config failed")
	}
	Config = iniConf
}
