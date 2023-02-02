package controllers

import (
	"ehang.io/nps/lib/file"
	"ehang.io/nps/server"
	"ehang.io/nps/server/tool"
	"encoding/json"
	"fmt"
	"github.com/beego/beego"
)

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Prepare() {
	if c.GetSession("clientId") == nil {
		c.CustomAbort(401, "未登录")
	}
}

func (c *AuthController) getClientId() int {
	clientId := c.GetSession("clientId")
	if clientId == nil {
		clientId = 1
	}

	return clientId.(int)
}

func (c *AuthController) Base() {
	data := make(map[string]any)

	data["user_account_url"] = c.GetSession("accountUrl")
	data["user_nickname"] = c.GetSession("nickname")
	data["user_avatar"] = c.GetSession("avatar")

	clientId := c.getClientId()
	client, err := file.GetDb().JsonDb.GetClient(clientId)
	if err != nil {
		c.CustomAbort(401, "未找到客户端 要不退出重新登录试试")
	}

	data["nps_vkey"] = client.VerifyKey
	data["nps_rate_limit"] = client.RateLimit
	data["nps_is_connect"] = client.IsConnect
	data["nps_flow_export"] = client.Flow.ExportFlow
	data["nps_flow_inlet"] = client.Flow.InletFlow

	data["nps_basic_username"] = client.Cnf.U
	data["nps_basic_password"] = client.Cnf.P

	data["nps_bridge_type"] = beego.AppConfig.String("BRIDGE_TYPE")
	data["nps_bridge_domain"] = beego.AppConfig.String("BRIDGE_DOMAIN")
	data["nps_bridge_port"] = beego.AppConfig.String("BRIDGE_PORT")

	data["nps_allow_ports"] = beego.AppConfig.String("ALLOW_PORTS")

	c.Data["json"] = data
	c.ServeJSON()
}

type Tunnel struct {
	Type   string
	Port   int
	Target string
	Remark string
}

func (c *AuthController) getTunnels() ([]*file.Tunnel, int) {
	return server.GetTunnel(0, 100, "", c.getClientId(), "")
}

func (c *AuthController) TunnelGet() {
	lists, count := c.getTunnels()

	var rows []map[string]any
	for i := 0; i < len(lists); i++ {
		list := lists[i]
		row := make(map[string]any)
		row["type"] = list.Mode
		row["port"] = list.Port
		row["target"] = list.Target.TargetStr
		row["remark"] = list.Remark
		rows = append(rows, row)
	}

	data := make(map[string]any)
	data["rows"] = rows
	data["count"] = count

	c.Data["json"] = data
	c.ServeJSON()
}

func (c *AuthController) TunnelPost() {
	var tunnels []Tunnel
	body := c.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &tunnels)
	if err != nil {
		println("json error", err.Error())
		c.CustomAbort(422, err.Error())
	}

	oldTunnels, _ := c.getTunnels()
	for i := 0; i < len(oldTunnels); i++ {
		_ = server.DelTask(oldTunnels[i].Id)
	}

	tunnelLens := len(tunnels)

	var client *file.Client

	if client, err = file.GetDb().GetClient(c.getClientId()); err != nil {
		c.CustomAbort(422, "客户端有误 请联系管理员")
	}

	if client.MaxTunnelNum != 0 && tunnelLens >= client.MaxTunnelNum {
		c.CustomAbort(422, fmt.Sprintf("目前最多只能有 %v 个隧道", client.MaxTunnelNum))
	}

	for i := 0; i < tunnelLens; i++ {
		tunnel := tunnels[i]

		id := int(file.GetDb().JsonDb.GetTaskId())
		t := &file.Tunnel{
			Port:      tunnel.Port,
			ServerIp:  "",
			Mode:      tunnel.Type,
			Target:    &file.Target{TargetStr: tunnel.Target, LocalProxy: false},
			Id:        id,
			Status:    true,
			Remark:    tunnel.Remark,
			Password:  "",
			LocalPath: "",
			StripPre:  "",
			Flow:      &file.Flow{},
		}

		if t.Port <= 0 {
			t.Port = tool.GenerateServerPort(t.Mode)
		} else {
			if !tool.TestServerPort(t.Port, t.Mode) {
				c.CustomAbort(422, fmt.Sprintf("服务端端口 %v 已被使用", t.Port))
			}
		}

		t.Client = client

		if err := file.GetDb().NewTask(t); err != nil {
			c.CustomAbort(422, err.Error())
		}
		if err := server.AddTask(t); err != nil {
			c.CustomAbort(422, err.Error())
		}
	}

	data := make(map[string]any)
	data["message"] = "保存成功"

	c.Data["json"] = data
	c.ServeJSON()
}
