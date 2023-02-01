package connection

import (
	"net"
	"os"
	"strconv"

	"github.com/beego/beego"
	"github.com/beego/beego/logs"
)

var bridgePort string
var httpsPort string
var httpPort string

func InitConnectionService() {
	bridgePort = beego.AppConfig.String("BRIDGE_PORT")
	httpsPort = beego.AppConfig.String("HTTPS_PROXY_PORT")
	httpPort = beego.AppConfig.String("HTTP_PROXY_PORT")
}

func GetBridgeListener(tp string) (net.Listener, error) {
	logs.Info("server start, the bridge type is %s, the bridge port is %s", tp, bridgePort)
	var p int
	var err error
	if p, err = strconv.Atoi(bridgePort); err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(beego.AppConfig.String("BRIDGE_IP")), p, ""})
}

func GetHttpListener() (net.Listener, error) {
	logs.Info("start http listener, port is", httpPort)
	return getTcpListener(beego.AppConfig.String("HTTP_PROXY_IP"), httpPort)
}

func GetHttpsListener() (net.Listener, error) {
	logs.Info("start https listener, port is", httpsPort)
	return getTcpListener(beego.AppConfig.String("HTTP_PROXY_IP"), httpsPort)
}

func getTcpListener(ip, p string) (net.Listener, error) {
	port, err := strconv.Atoi(p)
	if err != nil {
		logs.Error(err)
		os.Exit(0)
	}
	if ip == "" {
		ip = "0.0.0.0"
	}
	return net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
}
