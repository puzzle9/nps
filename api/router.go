package api

import (
	"ehang.io/nps/api/controllers"
	"github.com/beego/beego"
)

func Init() {
	beego.Router("/", &controllers.IndexController{}, "*:Start")
	beego.AutoRouter(&controllers.SignController{})
	beego.Router("/base", &controllers.AuthController{}, "post:Base")
	beego.Router("/updateVKey", &controllers.AuthController{}, "post:UpdateVKey")
	beego.Router("/tunnel", &controllers.AuthController{}, "get:TunnelLists")
	beego.Router("/tunnel", &controllers.AuthController{}, "post:TunnelCreate")
	beego.Router("/tunnel", &controllers.AuthController{}, "delete:TunnelDelete")
}
