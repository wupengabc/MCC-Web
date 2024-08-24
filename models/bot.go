package models

type Bots struct {
	ID         int64  `xorm:"'id' pk autoincr"`
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
