package controllers

import (
	"github.com/beego/beego"
	"time"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Start() {
	m := make(map[string]interface{})
	m["name"] = beego.BConfig.AppName
	m["time"] = time.Now().Unix()
	c.Data["json"] = m
	c.ServeJSON()
}
