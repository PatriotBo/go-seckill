package service

import "github.com/astaxie/beego"

type ApiController struct {
	beego.Controller
}

func InitServer() {
	serverPort := Config.String("port")
	beego.Router("/request", &ApiController{}, "post:HandleSecKill")
	beego.Run(serverPort)
}

func (ac *ApiController) HandleSecKill() {

}
