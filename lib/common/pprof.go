package common

import (
	"github.com/beego/beego"
	"github.com/beego/beego/logs"
	"net/http"
	_ "net/http/pprof"
)

func InitPProfFromFile() {
	ip := beego.AppConfig.String("PPROF_IP")
	p := beego.AppConfig.String("PPROF_PORT")
	if len(ip) > 0 && len(p) > 0 && IsPort(p) {
		runPProf(ip + ":" + p)
	}
}

func InitPProfFromArg(arg string) {
	if len(arg) > 0 {
		runPProf(arg)
	}
}

func runPProf(ipPort string) {
	go func() {
		_ = http.ListenAndServe(ipPort, nil)
	}()
	logs.Info("PProf debug listen on", ipPort)
}
