package models

import "fmt"

type Bots struct {
	I          int64  `xorm:"'id' pk autoincr"`
	Name       string `xorm:"'name' TEXT"`
	Version    string `xorm:"'version' TEXT default 'auto'"`
	Belong     string `xorm:"'belong' TEXT"`
	Open       string `xorm:"'open' TEXT"`
	Forge      string `xorm:"'forge' TEXT"`
	Server     string `xorm:"'server' TEXT"`
	Connection string `xorm:"'connection' TEXT"`
}

func InsertBot(name, version, belong, open, forge, server, connection string) {
	var bot Bots
	data := &Bots{
		Name:       name,
		Version:    version,
		Belong:     belong,
		Forge:      forge, // JSON 数字解析为 float64
		Open:       open,
		Server:     server,
		Connection: connection,
	}
	has, err := Db.Where("name = ? AND belong = ? AND server = ?", name, belong, server).Get(&bot)
	if err != nil {
		return
	}
	if has {
		_, err = Db.Where("name = ? AND belong = ? AND server = ?", name, belong, server).Update(data)
		return
	} else {
		_, err = Db.Insert(data)
		return
	}
}

func GetBotByUsername(username string) []Bots {
	var bots []Bots
	err := Db.Where("belong = ?", username).Find(&bots)
	if err != nil {
		return nil
	}
	return bots
}

type Commands struct {
	Id      int64  `xorm:"'id' INT(11) PK AUTOINCR"`
	Name    string `xorm:"'name' TEXT"`
	Belong  string `xorm:"'belong' TEXT"`
	Command string `xorm:"'command' TEXT"`
	Call    string `xorm:"'call' TEXT"`
}

func GetCommands(name, belong string) []Commands {
	var commands []Commands
	err := Db.Where("name = ? AND belong = ?", name, belong).Find(&commands)
	if err != nil {
		return nil
	}
	return commands
}

func AddCommand(name, belong, command, call string) string {
	data := Commands{
		Name:    name,
		Belong:  belong,
		Command: command,
		Call:    call,
	}
	has, err := Db.Where("name = ? AND belong = ? AND command = ? AND call = ?", name, belong, command, call).Get(&Commands{})
	if err != nil {
		return "2"
	}
	if has {
		return "1"
	}
	_, err = Db.Insert(data)
	return "0"
}

func DeleteCommand(name, belong, command, call string) string {
	_, err := Db.Where("name = ? AND belong = ? AND command = ? AND call = ?", name, belong, command, call).Delete(&Commands{})
	if err != nil {
		return "2"
	}
	return "0"
}

func ChangeCommand(name, belong, oldcommand, oldcall, newcommand, newcall string) string {
	data := Commands{
		Command: newcommand,
		Call:    newcall,
	}
	has, err := Db.Where("name = ? AND belong = ? AND command = ? AND call = ?", name, belong, oldcommand, oldcall).Get(&Commands{})
	if err != nil {
		return "2"
	}
	if has {
		_, err = Db.Where("name = ? AND belong = ? AND command = ? AND call = ?", name, belong, oldcommand, oldcall).Update(data)
		return "0"
	}
	return "1"
}

func GetCountByBelong(belong string) int64 {
	count, err := Db.Where("belong = ?", belong).Count(&Bots{})
	fmt.Println(count)
	if err != nil {
		return 100000
	}
	return count
}

func GetBotByNameBelongServer(name, belong, server string) Bots {
	var bot Bots
	_, err := Db.Where("name = ? AND belong = ? AND server = ?", name, belong, server).Get(&bot)
	if err != nil {
		return Bots{}
	}
	return bot
}

func DelBotByNameBelongServer(name, belong, server string) {
	_, err := Db.Where("name = ? AND belong = ? AND server = ?", name, belong, server).Delete(&Bots{})
	if err != nil {
		return
	}
	return
}
