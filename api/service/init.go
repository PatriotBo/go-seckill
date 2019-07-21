package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	log "github.com/astaxie/beego/logs"
)

type ApiController struct {
	beego.Controller
}

var (
	Config config.Configer
)

func Init() {
	initConfig()
	initRedis()
}

// 初始化redis
func initRedis() {

}

//初始化配置
func initConfig() {
	iniConf, err := config.NewConfig("ini", "/config/api.conf")
	if err != nil {
		log.Error("init config failed:", err)
		panic("init config failed")
		return
	}
	Config = iniConf
}
