package models

type Servers struct {
	Id         int64  `xorm:"pk autoincr 'id'"` // 主键，自增
	Server     string `xorm:"notnull 'server'"`
	Key        string `xorm:"notnull 'key'"`
	Permission int64  `xorm:"notnull default(0) 'permission'"`
	Name       string `xorm:"notnull 'name'"`
}

func GetAllServer() []Servers {
	var servers []Servers
	err := Db.Find(&servers)
	if err != nil {
		return nil
	}
	return servers
}

func GetServerByPermission(permission int64) []Servers {
	var servers []Servers
	err := Db.Where("permission <= ?", permission).Find(&servers)
	if err != nil {
		return nil
	}
	return servers
}

func GetServerKey(server string) string {
	var serverKey Servers
	_, err := Db.Where("server = ?", server).Get(&serverKey)
	if err != nil {
		return ""
	}
	return serverKey.Key
}

func GetServerById(id int64) string {
	var Server Servers
	_, err := Db.Where("id = ?", id).Get(&Server)
	if err != nil {
		return ""
	}
	return Server.Server
}
