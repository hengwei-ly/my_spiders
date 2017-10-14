package models

import (
	"fmt"
	"myMovies/utils"
	"path/filepath"

	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var Engine *xorm.Engine

func InitEngine() error {
	var err error

	params, err := utils.ReadProperties(filepath.Join("conf", "conf.properties"))
	if err != nil {
		return errors.New("读取初始化参数失败: " + err.Error())
	}

	driverName := params["driverName"]
	addr := params["addr"]
	port := params["port"]
	username := params["username"]
	password := params["password"]
	dbName := params["dbName"]

	Engine, err = xorm.NewEngine(driverName,
		fmt.Sprint(username, ":", password, "@tcp(", addr, ":", port, ")/", dbName, "?charset=utf8"))
	if err != nil {
		return errors.New("获取engine失败： " + err.Error())
	}
	return nil
}

func InitTable() error {
	myMovie := MyMovie{}
	flag, err := Engine.Exist(&myMovie)
	if err != nil || !flag {
		err = Engine.Sync(&myMovie)
		if err != nil {
			return err
		}
	}
	return nil
}
