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
	beego.Router("/getbotstatus", &controllers.MainController{}, "post:GetBotStatus")
	beego.Router("/manager", &controllers.MainController{}, "get:Manager")
	beego.Router("/ws/:param1/:param2", &controllers.WebSocketController{}, "get:ForwardWS")
	beego.Router("/startbot", &controllers.MainController{}, "post:StartBot")
	beego.Router("/stopbot", &controllers.MainController{}, "post:StopBot")
	beego.Router("/getbotlist", &controllers.MainController{}, "get:GetBotList")
	beego.Router("/getcommands", &controllers.MainController{}, "post:GetCommands")
	beego.Router("/getserverlist", &controllers.MainController{}, "get:GetServerList")
	beego.Router("/getadd", &controllers.MainController{}, "get:Getadd")
	beego.Router("/addbot", &controllers.MainController{}, "post:AddBot")
	beego.Router("/test", &controllers.MainController{}, "get:Test")
}
