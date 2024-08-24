package models

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var Db *xorm.Engine

func init() {
	var err error
	Db, err = xorm.NewEngine("sqlite3", "./data/database.db")
	if err != nil {
		fmt.Printf("不能连接到数据库: %v\n", err)
	}

	// 检查是否可以成功执行一个简单的查询
	if _, err = Db.Exec("SELECT 1"); err != nil {
		fmt.Printf("无法访问数据库: %v\n", err)
	} else {
		fmt.Println("数据库连接成功")
	}
	Db.SetMaxIdleConns(3)
	Db.SetMaxOpenConns(10)
	Db.ShowSQL(true)
}
