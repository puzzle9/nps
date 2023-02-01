package controllers

import (
	"github.com/beego/beego"
	"time"
)

type IndexController struct {
	beego.Controller
}

func (s *IndexController) Start() {
	m := make(map[string]interface{})
	m["name"] = beego.BConfig.AppName
	m["time"] = time.Now().Unix()
	s.Data["json"] = m
	s.ServeJSON()
}
