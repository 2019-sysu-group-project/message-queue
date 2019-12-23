package model

import (
	"fmt"
	"log"
	"os"

	"Mqservice/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var GormDB *gorm.DB

var maxConnectionTime = 5

func init() {
	fmt.Println("init函数1被执行！")
	times := 1
	for err := connectDB(); err != nil; times++ {
		if times == maxConnectionTime {
			panic(fmt.Sprint("can not connect to db after ", times, " times"))
			os.Exit(1)
			// break
		}
		log.Print("connect database with error", err, "reconnecting...")
	}
	// 将gorm调用的接口实时对应输出为真正执行的sql语句，用于debug使用
	GormDB.LogMode(true)

}

func reConnectDB() error {
	return connectDB()
}

func connectDB() error {
	db, err := gorm.Open(config.Mysql, config.Dbconnection+"?charset=utf8&parseTime=True") //这里的True首字母要大写！
	if err != nil {
		return err
	}
	//db.AutoMigrate(&User{}).AutoMigrate(&Product{}).AutoMigrate(&Service{})
	GormDB = db
	return nil
}
