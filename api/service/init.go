package service

import (
	"github.com/astaxie/beego/config"
	log "github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	"go-seckill/api/entity"
)

var (
	Config       config.Configer
	RedisClient  *entity.RedisConfig
	SecKillCtx   *entity.SecKillContext
	ConsulClient *api.Client
)

func Init() {
	initCtx()
	initConfig()
	initRedis()
	initConsul()
}

func initCtx() {
	SecKillCtx = &entity.SecKillContext{
		ProductInfoMap: make(map[int32]*entity.ProductInfo, 1024),
		RequestChan:    make(chan *entity.Request, 1024),
		UserConnMap:    make(map[string]chan *entity.Response, 1024),
	}
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
	RedisClient.Client = c
	RedisClient.IpBlackHash = Config.String("redis.blackIps")
	RedisClient.IdBlackHash = Config.String("redis.blackIds")
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

func initConsul() {
	addr := Config.String("consul.addr")
	conf := api.DefaultConfig()
	conf.Address = addr
	client, err := api.NewClient(conf)
	if err != nil {
		log.Error("init consul failed:", err)
		panic(err.Error())
	}
	ConsulClient = client
}
