package main

import (
	"new_project/models"
	_ "new_project/routers"

	"github.com/astaxie/beego"
)

func init() {
	models.RegisterMsql()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
