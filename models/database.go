package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser     = beego.AppConfig.String("dbUser")
	dbPassword = beego.AppConfig.String("dbPassword")
	dbHost     = beego.AppConfig.String("dbHost")
	dbPort     = beego.AppConfig.String("dbPort")
	dbName     = beego.AppConfig.String("dbName")
)

func RegisterMsql() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPassword, dbHost, dbPort, dbName)
	orm.RegisterModel(new(AdvRecord), new(UploadRecord))
	orm.RegisterDataBase("default", "mysql", dataSource, 30)
}
