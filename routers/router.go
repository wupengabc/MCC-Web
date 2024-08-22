package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"mcc_web/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
