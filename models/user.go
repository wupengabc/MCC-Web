package models

type Users struct {
	Id         int64  `xorm:"pk autoincr 'id'"`                // 主键，自增
	Username   string `xorm:"notnull 'username'"`              // 用户名，不能为空
	Password   string `xorm:"notnull 'password'"`              // 密码，不能为空
	Permission int    `xorm:"notnull default(0) 'permission'"` // 权限，默认值为0
}

type Regcodes struct {
	Id   int64  `xorm:"pk autoincr 'id'"` // 主键，自增
	Code string `xorm:"notnull 'code'"`
}

func InsertCode(code string) {
	var codes Regcodes
	codes.Code = code
	_, err := Db.Insert(&codes)
	if err != nil {
		return
	}
	return
}

func CheckCode(code string) int {
	var codes Regcodes
	result, err := Db.Where("code = ?", code).Get(&codes)
	if err != nil {
		return 2
	}
	if result {
		_, err := Db.Where("code = ?", code).Delete(&codes)
		if err != nil {
			return 2
		}
		return 1
	}
	return 0
}

func InsertUser(username, password string) int {
	var user Users
	user.Username = username
	user.Password = password
	user.Permission = 1
	status, err2 := Db.Where("username = ?", username).Get(&user)
	if err2 != nil {
		return 0
	}
	if status {
		return 2
	}
	_, err := Db.Insert(&user)
	if err != nil {
		return 0
	}
	return 1
}

func CheckUser(username, password string) (int, error, int) {
	var user Users
	result, err := Db.Where("username = ? and password = ?", username, password).Get(&user)
	status, err2 := Db.Where("username = ?", username).Get(&user)
	if err != nil {
		return 0, err, 0
	}
	if err2 != nil {
		return 0, err2, 0
	}
	if status {
		if result {
			permission := user.Permission
			return 1, nil, permission
		} else {
			return 2, nil, 0
		}
	} else {
		return 3, nil, 0
	}
}
