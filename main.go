package main

import (
	_"phantom/routers"
	"github.com/astaxie/beego"
	"phantom/tools"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.AddFuncMap("addone",tools.Addone)
	beego.Run()
}
