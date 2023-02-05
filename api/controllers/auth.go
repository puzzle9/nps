package controllers

import (
	"ehang.io/nps/api/service"
	"ehang.io/nps/lib/crypt"
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

func (c *AuthController) UpdateVKey() {
	client, err := file.GetDb().JsonDb.GetClient(c.getClientId())
	if err != nil {
		c.CustomAbort(422, "客户端 异常")
	}

	if client.IsConnect {
		c.CustomAbort(422, "将客户端断开后再进行更新")
	}

	clientId := client.Id

	server.DelTunnelAndHostByClientId(clientId, false)
	server.DelClientConnect(clientId)

	vkey := crypt.GetRandomString(16)
	client.VerifyKey = vkey

	file.GetDb().JsonDb.StoreClientsToJsonFile()

	data := make(map[string]any)
	data["nps_vkey"] = vkey
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

func (c *AuthController) TunnelLists() {
	lists, count := c.getTunnels()

	var rows []map[string]any
	for i := 0; i < len(lists); i++ {
		list := lists[i]
		rowId, _ := service.HashIdEncode(list.Id)

		var urlDomain string
		hostId := list.HostId

		if hostId != 0 {
			if host, err := file.GetDb().GetHostById(hostId); err == nil {
				urlDomain = host.Host
			}
		}

		row := make(map[string]any)
		row["id"] = rowId
		row["type"] = list.Mode
		row["port"] = list.Port
		row["target"] = list.Target.TargetStr
		row["remark"] = list.Remark
		row["url_domain"] = urlDomain
		row["url_ip"] = fmt.Sprintf("%v:%v", beego.AppConfig.String("BRIDGE_DOMAIN"), list.Port)
		rows = append(rows, row)
	}

	data := make(map[string]any)
	data["rows"] = rows
	data["count"] = count

	c.Data["json"] = data
	c.ServeJSON()
}

func (c *AuthController) TunnelCreate() {
	var tunnel Tunnel
	body := c.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &tunnel)
	if err != nil {
		println("json error", err.Error())
		c.CustomAbort(422, err.Error())
	}

	var client *file.Client

	if client, err = file.GetDb().GetClient(c.getClientId()); err != nil {
		c.CustomAbort(422, "客户端有误 请联系管理员")
	}

	if client.MaxTunnelNum != 0 && client.GetTunnelNum() >= client.MaxTunnelNum {
		c.CustomAbort(422, fmt.Sprintf("目前最多只能有 %v 个隧道", client.MaxTunnelNum))
	}

	tunnelPort := tunnel.Port
	tunnelType := tunnel.Type

	switch tunnelType {
	case "tcp":
	case "udp":
	case "socks5":
	case "httpProxy":
		tunnelType = tunnelType
		break
	default:
		c.CustomAbort(422, "暂不支持此模式")
		break
	}

	target := &file.Target{TargetStr: tunnel.Target, LocalProxy: false}

	if tunnelPort <= 0 {
		tunnelPort = tool.GenerateServerPort(tunnelType)
	} else {
		if !tool.TestServerPort(tunnelPort, tunnelType) {
			c.CustomAbort(422, fmt.Sprintf("服务端端口 %v 已被使用", tunnelPort))
		}
	}

	var hostId int

	if tunnelType == "tcp" {
		hostId = int(file.GetDb().JsonDb.GetHostId())
		h := &file.Host{
			Id:           hostId,
			Host:         fmt.Sprintf("%v-%v.%v", tunnelPort, tunnelType, beego.AppConfig.String("HTTP_PROXY_HOST")),
			Target:       target,
			HeaderChange: "",
			HostChange:   "",
			Remark:       "",
			Location:     "",
			Flow:         &file.Flow{},
			Scheme:       "http",
		}

		h.Client = client

		if err := file.GetDb().NewHost(h); err != nil {
			c.CustomAbort(422, err.Error())
		}
	}

	tunnelId := int(file.GetDb().JsonDb.GetTaskId())
	t := &file.Tunnel{
		Port:      tunnelPort,
		ServerIp:  "",
		Mode:      tunnelType,
		Target:    target,
		Id:        tunnelId,
		Status:    true,
		Remark:    tunnel.Remark,
		Password:  "",
		LocalPath: "",
		StripPre:  "",
		Flow:      &file.Flow{},
		HostId:    hostId,
	}

	t.Client = client

	if err := file.GetDb().NewTask(t); err != nil {
		c.CustomAbort(422, err.Error())
	}
	if err := server.AddTask(t); err != nil {
		c.CustomAbort(422, err.Error())
	}

	data := make(map[string]any)
	data["message"] = "添加成功"

	c.Data["json"] = data
	c.ServeJSON()
}

func (c *AuthController) TunnelDelete() {
	ids, err := service.HashIdDecode(c.GetString("id"))
	if err != nil {
		c.CustomAbort(422, err.Error())
	}

	tunnel, err := file.GetDb().GetTask(ids[0])
	if err != nil {
		c.CustomAbort(422, err.Error())
	}

	_ = file.GetDb().DelHost(tunnel.HostId)
	_ = server.DelTask(ids[0])

	data := make(map[string]any)
	data["message"] = "删除成功"

	c.Data["json"] = data
	c.ServeJSON()
}
