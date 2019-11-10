package routers

import (
	"github.com/astaxie/beego"
	"new_project/controllers"
)

func init() {
	beego.Router("/upload", &controllers.UploadController{})
}
