package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDB *gorm.DB

func init(){
	MysqlDB,err:=gorm.Open("mysql", "lcx:123456@tcp(127.0.0.1:3306)/mFPS?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	defer MysqlDB.Close()
}

