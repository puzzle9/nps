package api

import (
	"ehang.io/nps/api/controllers"
	"github.com/beego/beego"
)

func Init() {
	beego.Router("/", &controllers.IndexController{}, "*:Start")
}
