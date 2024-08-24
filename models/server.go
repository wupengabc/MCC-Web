package models

type Servers struct {
	Id         int64  `xorm:"pk autoincr 'id'"` // 主键，自增
	Server     string `xorm:"notnull 'server'"`
	Key        string `xorm:"notnull 'key'"`
	Permission int64  `xorm:"notnull default(0) 'permission'"`
}

func GetAllServer() []Servers {
	var servers []Servers
	err := Db.Find(&servers)
	if err != nil {
		return nil
	}
	return servers
}

func GetServerByPermission(permission string) []Servers {
	var servers []Servers
	err := Db.Where("permission <= ?", permission).Find(&servers)
	if err != nil {
		return nil
	}
	return servers
}
