package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"mcc_web/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Index() {
	c.TplName = "index.html"
	return
}

func (c *MainController) LoginGet() {
	if c.GetSession("username") != nil {
		c.TplName = "login.html"
		c.Data["IsLogin"] = 1
	}
	c.TplName = "login.html"
	return
}

func (c *MainController) LoginPost() {
	username := c.GetString("username")
	password := c.GetString("password")
	result, err, permission := models.CheckUser(username, password)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "数据文件错误，请检查后台报错"}
		err := c.ServeJSON()
		if err != nil {
			log.Fatalf("Error: %v", err)
			return
		}
	} else {
		if result == 0 {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "数据文件错误，请检查后台报错"}
			err := c.ServeJSON()
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		} else if result == 1 {
			c.Data["json"] = map[string]interface{}{"code": 200, "msg": "登录成功"}
			err := c.SetSession("username", username)
			if err != nil {
				return
			}
			err = c.SetSession("permission", permission)
			if err != nil {
				return
			}
			err = c.ServeJSON()
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		} else if result == 2 {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "用户名或密码错误"}
			err := c.ServeJSON()
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		} else if result == 3 {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "用户不存在"}
			err := c.ServeJSON()
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "数据文件错误，请检查后台报错"}
			err := c.ServeJSON()
			if err != nil {
				log.Fatalf("Error: %v", err)
				return
			}
		}
	}
}

func (c *MainController) RegGet() {
	if c.GetSession("username") != nil {
		c.TplName = "login.html"
		c.Data["IsLogin"] = 1
	}
	c.TplName = "reg.html"
	return
}

func (c *MainController) RegPost() {
	username := c.GetString("username")
	password := c.GetString("password")
	code := c.GetString("code")
	CodeStatus := models.CheckCode(code)
	if CodeStatus == 1 {
		result := models.InsertUser(username, password)
		if result == 1 {
			c.Data["json"] = map[string]interface{}{"code": 200, "msg": "注册成功,即将为你跳转登录界面"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else if result == 2 {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "用户已存在,请更换用户名后再试"}
			models.InsertCode(code)
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else if result == 0 {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "注册失败,请重新注册或联系管理员检查后台报错"}
			models.InsertCode(code)
			err := c.ServeJSON()
			if err != nil {
				return
			}
		}
	} else if CodeStatus == 0 {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "注册码错误,请重新输入"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	} else if CodeStatus == 2 {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "注册码验证出错,请重新注册或联系管理员检查后台报错"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

func (c *MainController) Panel() {
	username := c.GetSession("username")
	permission := c.GetSession("permission")
	if c.GetSession("username") != nil {
		c.TplName = "panel.html"
		c.Data["username"] = username
		c.Data["permission"] = permission
		return
	}
	c.Redirect("/login", 302)
	return
}

func (c *MainController) Logout() {
	err := c.DestroySession()
	if err != nil {
		return
	}
	c.Redirect("/", 302)
}
