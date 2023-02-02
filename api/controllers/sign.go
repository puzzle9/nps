package controllers

import (
	"ehang.io/nps/api/service"
	"ehang.io/nps/lib/file"
	"github.com/beego/beego"
)

type SignController struct {
	beego.Controller
}

func (c *SignController) In() {
	c.Redirect(service.GetSignInUrl(), 302)
}

func (c *SignController) Out() {
	c.DestroySession()
	c.Redirect("/", 302)
}

func (c *SignController) Oauth() {
	token := c.GetString("token")
	body, err := service.GetUserInfoByToken(token)
	if err != nil {
		c.CustomAbort(500, "账户 登录失败")
		return
	}

	encodeUserId := body.EncodeUserId

	clientId := 0
	file.GetDb().JsonDb.Clients.Range(func(key, value any) bool {
		client := value.(*file.Client)
		if client.EncodeUserId == encodeUserId {
			clientId = client.Id
			return false
		}
		return true
	})

	if clientId == 0 {
		cnf := &file.Config{
			Compress: true,
			Crypt:    true,
		}

		clientId = int(file.GetDb().JsonDb.GetClientId())

		client := &file.Client{
			Id:              clientId,
			EncodeUserId:    encodeUserId,
			Status:          true,
			Cnf:             cnf,
			ConfigConnAllow: false,
			RateLimit:       500,
			Flow:            &file.Flow{},
			MaxTunnelNum:    10,
		}
		if err := file.GetDb().NewClient(client); err != nil {
			c.CustomAbort(500, "账户 注册失败")
		}
	}

	c.SetSession("clientId", clientId)
	c.SetSession("accountUrl", body.AccountUrl)
	c.SetSession("encodeUserId", encodeUserId)
	c.SetSession("nickname", body.Nickname)
	c.SetSession("avatar", body.Avatar)

	c.Redirect(beego.AppConfig.String("WEB_URL_HOME"), 302)
}
