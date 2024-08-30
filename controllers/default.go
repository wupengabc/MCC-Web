package controllers

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
	"log"
	"mcc_web/models"
	"net/http"
	"strconv"
)

type MainController struct {
	beego.Controller
}

type WebSocketController struct {
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

func (c *MainController) GetBotList() {
	username := c.GetSession("username")
	if username != nil {
		data := models.GetBotByUsername(username.(string))
		c.Data["json"] = data
		err := c.ServeJSON()
		if err != nil {
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

type GetBotStatus struct {
	Name   string `json:"name"`
	Server string `json:"server"`
	Belong string `json:"belong"`
}

func (c *MainController) GetBotStatus() {
	if c.GetSession("username") != nil {
		body := c.Ctx.Input.RequestBody
		var data GetBotStatus
		if err := json.Unmarshal(body, &data); err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Data["json"] = map[string]string{"error": "Invalid JSON"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
		id, _ := strconv.Atoi(data.Server)
		server := models.GetServerById(int64(id))
		key := models.GetServerKey(server)
		url := "http://" + server + "/" + key + "/status"
		datatostr, err := json.Marshal(data)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.Data["json"] = map[string]string{"error": "Invalid JSON"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
		status := models.PostData(url, string(datatostr))
		c.Data["json"] = map[string]interface{}{"code": 200, "status": status}
		err = c.ServeJSON()
		if err != nil {
			return
		}
	}
}

func (c *MainController) Logout() {
	err := c.DestroySession()
	if err != nil {
		return
	}
	c.Redirect("/", 302)
}

func (c *MainController) Manager() {
	if c.GetSession("username") != nil {
		c.TplName = "manager.html"
		return
	} else {
		c.Redirect("/", 302)
		return
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *WebSocketController) ForwardWS() {
	param2 := c.Ctx.Input.Param(":param2")
	username := fmt.Sprintf("%v", c.GetSession("username"))
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Redirect("/", 302)
		return
	}
	if username != param2 && permission < 7 {
		c.Redirect("/", 302)
		return
	}
	if c.GetSession("username") != nil {
		// 获取请求路径中的参数
		param1 := c.Ctx.Input.Param(":param1")
		id, _ := strconv.Atoi(c.GetString("server"))

		// 获取目标 WebSocket 地址
		server := models.GetServerById(int64(id))
		key := models.GetServerKey(server)
		targetWSURL := "ws://" + server + "/" + key + "/" + param2 + "/" + param1
		fmt.Println("Connecting to target WebSocket URL:", targetWSURL)

		// 升级为 WebSocket 连接
		clientConn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
		if err != nil {
			fmt.Println("Failed to upgrade WebSocket connection:", err)
			return
		}
		defer clientConn.Close()

		// 连接目标 WebSocket 服务器
		serverConn, _, err := websocket.DefaultDialer.Dial(targetWSURL, nil)
		if err != nil {
			fmt.Println("Failed to connect to target WebSocket server:", err)
			return
		}
		defer func(serverConn *websocket.Conn) {
			err := serverConn.Close()
			if err != nil {

			}
		}(serverConn)

		// 转发消息
		go func() {
			for {
				_, message, err := clientConn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading from client WebSocket:", err)
					return
				}
				if err := serverConn.WriteMessage(websocket.TextMessage, message); err != nil {
					fmt.Println("Error writing to server WebSocket:", err)
					return
				}
			}
		}()

		for {
			_, message, err := serverConn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading from server WebSocket:", err)
				return
			}
			if err := clientConn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("Error writing to client WebSocket:", err)
				return
			}
		}
	} else {
		c.Redirect("/", 302)
		return
	}
}

type SwitchBots struct {
	Name   string `json:"name"`
	Belong string `json:"belong"`
	Server string `json:"server"`
}

func (c *MainController) StartBot() {
	username := c.GetSession("username").(string)
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	body := c.Ctx.Input.RequestBody
	var data SwitchBots
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if data.Belong == username || permission >= 7 {
		id, _ := strconv.Atoi(data.Server)
		server := models.GetServerById(int64(id))
		key := models.GetServerKey(server)
		status_url := "http://" + server + "/" + key + "/status"
		data2 := map[string]interface{}{"name": data.Name, "belong": data.Belong}
		jsontostr2, _ := json.Marshal(data2)
		result2 := models.PostData(status_url, string(jsontostr2))
		if result2 == "1" {
			c.Data["json"] = map[string]interface{}{"code": 401, "message": "机器人正在运行中,请勿重复开启"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else if result2 == "2" || result2 == "0" {
			url := "http://" + server + "/" + key + "/start"
			json1 := map[string]interface{}{"name": data.Name, "belong": data.Belong}
			jsontostr, _ := json.Marshal(json1)
			result := models.PostData(url, string(jsontostr))
			if result == "1" {
				c.Data["json"] = map[string]interface{}{"code": 200, "message": "成功发送开启请求"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			} else {
				c.Data["json"] = map[string]interface{}{"code": 401, "message": "请求失败，请检查服务器状态"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			}
		} else if result2 == "3" {
			c.Data["json"] = map[string]interface{}{"code": 401, "message": "后端出错, 无法开启机器人, 请联系管理员处理"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "message": "请求失败，请稍后重试"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "message": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

func (c *MainController) StopBot() {
	username := c.GetSession("username").(string)
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	body := c.Ctx.Input.RequestBody
	var data SwitchBots
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if data.Belong == username || permission >= 7 {
		id, _ := strconv.Atoi(data.Server)
		server := models.GetServerById(int64(id))
		key := models.GetServerKey(server)
		url := "http://" + server + "/" + key + "/stop"
		json1 := map[string]interface{}{"name": data.Name, "belong": data.Belong}
		jsontostr, _ := json.Marshal(json1)
		result := models.PostData(url, string(jsontostr))
		if result == "1" {
			c.Data["json"] = map[string]interface{}{"code": 200, "message": "成功发送关闭请求"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "message": "请求失败，请检查服务器状态"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "message": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

type GetCommands struct {
	Name   string `json:"name"`
	Belong string `json:"belong"`
}

func (c *MainController) GetCommands() {
	username := c.GetSession("username").(string)
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	body := c.Ctx.Input.RequestBody
	var data GetCommands
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if data.Belong == username || permission >= 7 {
		result := models.GetCommands(data.Name, data.Belong)
		c.Data["json"] = result
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "message": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

func (c *MainController) GetServerList() {
	username := c.GetSession("username")
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	if username != nil {
		servers := models.GetServerByPermission(int64(permission))
		var idlist []int64
		var namelist []string
		for _, item := range servers {
			id := item.Id
			name := item.Name
			idlist = append(idlist, id)
			namelist = append(namelist, name)
		}
		c.Data["json"] = map[string]interface{}{"code": 200, "data": map[string]interface{}{"idlist": idlist, "namelist": namelist}}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "message": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

func (c *MainController) Getadd() {
	username := c.GetSession("username")
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	if username != nil {
		has := models.GetCountByBelong(username.(string))
		num := int64(permission)
		if has < num {
			content := strconv.FormatInt(has, 10) + "/" + strconv.FormatInt(num, 10)
			c.Data["json"] = map[string]interface{}{"code": 200, "data": content}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else {
			content := "你的机器人已到达上限, 请升级权限或者删除旧的机器人, 你的机器人数量: " + strconv.FormatInt(has, 10) + "/" + strconv.FormatInt(num, 10)
			c.Data["json"] = map[string]interface{}{"code": 401, "data": content}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "message": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

type AddBot struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Forge      string `json:"forge"`
	Connection string `json:"connection"`
	Server     string `json:"server"`
}

func (c *MainController) AddBot() {
	username := c.GetSession("username")
	permission, err3 := strconv.Atoi(fmt.Sprintf("%v", c.GetSession("permission")))
	if err3 != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
	body := c.Ctx.Input.RequestBody
	var data AddBot
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if username != nil {
		has := models.GetCountByBelong(username.(string))
		num := int64(permission)
		if has < num {
			data3 := map[string]interface{}{"name": data.Name, "version": data.Version, "forge": data.Forge, "connection": data.Connection, "server": data.Server, "belong": username.(string)}
			data3tojson, _ := json.Marshal(data3)
			id, _ := strconv.Atoi(data.Server)
			server := models.GetServerById(int64(id))
			key := models.GetServerKey(server)
			url := "http://" + server + "/" + key + "/addbot"
			result := models.PostData(url, string(data3tojson))
			if result == "0" {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "添加失败, 当前用户下已存在同名机器人"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			} else if result == "1" {
				c.Data["json"] = map[string]interface{}{"code": 200, "msg": "添加成功, 请等待数据同步, 约5s"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			} else if result == "2" {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "添加失败, 服务端出错"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			} else {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "添加失败, 未知错误"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			}
		} else {
			content := "你的机器人已到达上限, 请升级权限或者删除旧的机器人, 你的机器人数量: " + strconv.FormatInt(has, 10) + "/" + strconv.FormatInt(num, 10)
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": content}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

func (c *MainController) Test() {
	result := models.GetCommands("7728", "admin")
	c.Data["json"] = result
	err := c.ServeJSON()
	if err != nil {
		return
	}
}
