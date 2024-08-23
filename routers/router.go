package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"mcc_web/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Index")
	beego.Router("/login", &controllers.MainController{}, "get:LoginGet;post:LoginPost")
	beego.Router("/reg", &controllers.MainController{}, "get:RegGet;post:RegPost")
	beego.Router("/panel", &controllers.MainController{}, "get:Panel")
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")
}
