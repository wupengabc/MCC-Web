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
	if username == "" || password == "" {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "用户名或密码不能为空"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
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
		if username == "" || password == "" {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "用户名或密码不能为空"}
			models.InsertCode(code)
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
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
		defer func(clientConn *websocket.Conn) {
			err := clientConn.Close()
			if err != nil {
			}
		}(clientConn)

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
		c.Data["json"] = map[string]interface{}{"code": 200, "data": result}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

type AddDeleteCommand struct {
	Name    string `json:"name"`
	Belong  string `json:"belong"`
	Command string `json:"command"`
	Call    string `json:"call"`
	Method  string `json:"method"`
}

func (c *MainController) AddDeleteCommand() {
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
	var data AddDeleteCommand
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if username == data.Belong || permission >= 7 {
		if data.Method == "add" {
			result := models.AddCommand(data.Name, data.Belong, data.Command, data.Call)
			c.Data["json"] = map[string]interface{}{"code": 200, "data": result}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else if data.Method == "delete" {
			result := models.DeleteCommand(data.Name, data.Belong, data.Command, data.Call)
			c.Data["json"] = map[string]interface{}{"code": 200, "data": result}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "data": "错误的Method"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
}

type ChangeCommand struct {
	Name       string `json:"name"`
	Belong     string `json:"belong"`
	Oldcommand string `json:"oldcommand"`
	Oldcall    string `json:"oldcall"`
	Newcommand string `json:"newcommand"`
	Newcall    string `json:"newcall"`
}

func (c *MainController) ChangeCommand() {
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
	var data ChangeCommand
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if username == data.Belong || permission >= 7 {
		result := models.ChangeCommand(data.Name, data.Belong, data.Oldcommand, data.Oldcall, data.Newcommand, data.Newcall)
		c.Data["json"] = map[string]interface{}{"code": 200, "data": result}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
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
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
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
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
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
			result4 := models.GetBotByNameBelongServer(data.Name, username.(string), data.Server)
			if result4 != (models.Bots{}) {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "添加失败, 当前用户下已存在同名机器人"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			}
			data3 := map[string]interface{}{"name": data.Name, "version": data.Version, "forge": data.Forge, "connection": data.Connection, "server": data.Server, "belong": username.(string)}
			data3tojson, _ := json.Marshal(data3)
			id, _ := strconv.Atoi(data.Server)
			server := models.GetServerById(int64(id))
			key := models.GetServerKey(server)
			fmt.Println(server)
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
				models.InsertBot(data.Name, data.Version, username.(string), "0", data.Forge, data.Server, data.Connection)
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

func (c *MainController) GetBotConfig() {
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
	if username == data.Belong || permission >= 7 {
		config := models.GetBotByNameBelongServer(data.Name, data.Belong, data.Server)
		if config != (models.Bots{}) {
			c.Data["json"] = map[string]interface{}{"code": 200, "data": config}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "data": "配置获取失败, 配置未同步或者网络错误"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

type UpdateBot struct {
	Oldname       string `json:"oldname"`
	Oldbelong     string `json:"oldbelong"`
	Oldserver     string `json:"oldserver"`
	Newname       string `json:"newname"`
	Newbelong     string `json:"newbelong"`
	Newversion    string `json:"newversion"`
	Newforge      string `json:"newforge"`
	Newconnection string `json:"newconnection"`
}

func (c *MainController) UpdateBot() {
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
	var data UpdateBot
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if username == data.Oldbelong || permission >= 7 {
		if username != data.Newbelong {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "你没有权限修改机器人所有者"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		}
		newconfig := map[string]interface{}{"oldname": data.Oldname, "oldbelong": username, "newname": data.Newname, "newbelong": data.Newbelong, "newversion": data.Newversion, "newforge": data.Newforge, "newconnection": data.Newconnection}
		configtostr, _ := json.Marshal(newconfig)
		id, _ := strconv.Atoi(data.Oldserver)
		server := models.GetServerById(int64(id))
		key := models.GetServerKey(server)
		url := "http://" + server + "/" + key + "/updatebot"
		models.DelBotByNameBelongServer(data.Oldname, username.(string), data.Oldserver)
		result := models.PostData(url, string(configtostr))
		if result == "0" {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "修改失败, 该机器人不存在"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else if result == "1" {
			c.Data["json"] = map[string]interface{}{"code": 200, "msg": "修改成功, 请等待数据同步后手动重启Bot, 同步时间约5s"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else if result == "2" {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "修改失败, 服务端错误, 请重试或者联系管理员处理"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else if result == "3" {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "同一服务器下, 您只能存在一个相同名称的Bot, 请尝试更换其他名字或者删除冲突的Bot, 或使用其他服务器"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "获取修改状态失败, 请5s后刷新页面查看是否修改成功或者联系管理员处理"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

func (c *MainController) DeleteBot() {
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
	if username == data.Belong || permission >= 7 {
		deletedata := map[string]interface{}{"name": data.Name, "belong": data.Belong}
		deletedatatostr, _ := json.Marshal(deletedata)
		id, _ := strconv.Atoi(data.Server)
		server := models.GetServerById(int64(id))
		key := models.GetServerKey(server)
		url := "http://" + server + "/" + key + "/removebot"
		models.DelBotByNameBelongServer(data.Name, data.Belong, data.Server)
		result := models.PostData(url, string(deletedatatostr))
		if result == "1" {
			c.Data["json"] = map[string]interface{}{"code": 200, "msg": "删除成功, 请等待数据同步, 约5s"}
			err := c.ServeJSON()
			models.DelBotByNameBelongServer(data.Name, data.Belong, data.Server)
			if err != nil {
				return
			}
		} else {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "获取删除结果失败, 请5s后刷新页面查看是否删除"}
			err := c.ServeJSON()
			models.DelBotByNameBelongServer(data.Name, data.Belong, data.Server)
			if err != nil {
				return
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "msg": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

func (c *MainController) GetNotice() {
	username := c.GetSession("username")
	if username != nil {
		notice := models.GetNotices()
		c.Data["json"] = map[string]interface{}{"code": 200, "data": notice}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

func (c *MainController) Setting() {
	username := c.GetSession("username")
	permission := c.GetSession("permission")
	if c.GetSession("username") != nil {
		c.TplName = "setting.html"
		c.Data["username"] = username
		c.Data["permission"] = permission
		count := models.GetCountByBelong(username.(string))
		c.Data["count"] = count
		return
	}
	c.Redirect("/login", 302)
	return
}

type ChangeUsername struct {
	Oldusername string `json:"oldusername"`
	Newusername string `json:"newusername"`
}

func (c *MainController) ChangeUsername() {
	username := c.GetSession("username")
	permission := c.GetSession("permission")
	intpermission, _ := strconv.Atoi(fmt.Sprintf("%v", permission))
	body := c.Ctx.Input.RequestBody
	var data ChangeUsername
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if username.(string) == data.Oldusername || intpermission > 7 {
		count := models.GetCountByBelong(username.(string))
		if count > 0 {
			c.Data["json"] = map[string]interface{}{"code": 401, "data": "请删除全部机器人后再尝试修改"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
			return
		} else {
			status := models.GetUser(data.Newusername)
			if status {
				c.Data["json"] = map[string]interface{}{"code": 401, "data": "用户名已存在"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
				return
			} else {
				if data.Newusername == "" {
					c.Data["json"] = map[string]interface{}{"code": 401, "data": "新用户名不能为空"}
					err := c.ServeJSON()
					if err != nil {
						return
					}
					return
				} else {
					result := models.ChangeUsername(data.Oldusername, data.Newusername)
					if result {
						c.Data["json"] = map[string]interface{}{"code": 200, "data": "修改成功"}
						err := c.DestroySession()
						if err != nil {
							return
						}
						err = c.ServeJSON()
						if err != nil {
							return
						}
						return
					} else {
						c.Data["json"] = map[string]interface{}{"code": 401, "data": "修改失败, 请重试"}
						err := c.ServeJSON()
						if err != nil {
							return
						}
						return
					}
				}
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足, 你无法修改其他人的用户名"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
	}
}

type ChangePassword struct {
	Username    string `json:"username"`
	Oldpassword string `json:"oldpassword"`
	Newpassword string `json:"newpassword"`
}

func (c *MainController) ChangePassword() {
	username := c.GetSession("username")
	permission := c.GetSession("permission")
	intpermission, _ := strconv.Atoi(fmt.Sprintf("%v", permission))
	body := c.Ctx.Input.RequestBody
	var data ChangePassword
	if err := json.Unmarshal(body, &data); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
		return
	}
	if username.(string) == data.Username || intpermission >= 7 {
		if data.Newpassword == "" {
			c.Data["json"] = map[string]interface{}{"code": 401, "msg": "密码不能为空"}
			err := c.ServeJSON()
			if err != nil {
				return
			}
		} else {
			result := models.ChangePassword(data.Username, data.Oldpassword, data.Newpassword)
			if result == 0 {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "服务器错误, 请稍后再试"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
			} else if result == 1 {
				c.Data["json"] = map[string]interface{}{"code": 200, "msg": "密码修改成功"}
				err := c.DestroySession()
				if err != nil {
					return
				}
				err = c.ServeJSON()
				if err != nil {
					return
				}
			} else if result == 2 {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "旧密码错误"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
			} else {
				c.Data["json"] = map[string]interface{}{"code": 401, "msg": "未知错误, 请稍后再试"}
				err := c.ServeJSON()
				if err != nil {
					return
				}
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"code": 401, "data": "权限不足"}
		err := c.ServeJSON()
		if err != nil {
			return
		}
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
