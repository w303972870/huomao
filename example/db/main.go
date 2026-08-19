package main

import (
	"fmt"
	"github.com/w303972870/huomao/db"
)

func main() {
	db := huomao_db.GetDbClass("sqlite3")
	db.Connect("/Users/huomao/Share/Golang/captcha/keys/sqlite3.db")
	const create_talbe_sql string = "CREATE TABLE IF NOT EXISTS `test_table` (" +
		"`cn` varchar(64) PRIMARY KEY," +
		"`article` text NOT NULL)"
	db.ExecSql(create_talbe_sql)
	i, _ := db.Insert(fmt.Sprint("insert into test_table values(", "\"huomao\"", ",\"laiye\"", ")"))
	fmt.Println("数据库插入成功：", i)
}
