package main

import (
	"errors"
	"flag"
	"fmt"
	"myMovies/models"
	"myMovies/services"
)

var modules map[string]models.ModuleInfo
var params map[string]string

func init() {
	err := models.InitEngine()
	if err != nil {
		fmt.Println("初始化engine失败")
		panic(err)
	}

	err = models.InitTable()
	if err != nil {
		fmt.Println("初始化table失败")
		panic(err)
	}
	modules, err = models.InitModuleInfos()
	if err != nil {
		fmt.Println("初始化model信息失败")
		panic(err)
	}

	initParams()
}

//初始化命令行参数
func initParams() {
	params = map[string]string{}
	updateTableCatelog := flag.String("m", "", "更新当前数据库中的信息, 参数为模块代号,-m all 为更新所有模块")

	updateStatusCatelog := flag.String("u", "", "下载完成后更新数据库中地址的状态,参数为模块代号，all 表示全部地址状态的更新")
	flag.Parse()

	params["updateTableCatelog"] = *updateTableCatelog
	params["updateStatusCatelog"] = *updateStatusCatelog

}

func main() {
	updateTableCatelog := params["updateTableCatelog"]
	updateStatusCatelog := params["updateStatusCatelog"]
	if len(updateTableCatelog) > 0 {
		err := updateTableByParam(updateTableCatelog)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	} else if len(updateStatusCatelog) > 0 {
		fmt.Println("更新状态看看" + updateStatusCatelog)
	} else {

		fmt.Println("以下为模块对应代码: ")
		for mun, mod := range modules {
			fmt.Println("\t" + mun + " ---> " + mod.Name)
		}
		fmt.Println()
		fmt.Println(`请输入参数:
	-m 更新数据库中的数据信息，拉取目标网页最新信息;
	-u 更新所有记录的状态，扫描目录，已下载的文件会被标识，同时更新`, models.DirPath, `下对应模块的.txt文件;`)
	}
}

//根据参数更新所有的url连接
func updateTableByParam(updateTableCatelog string) error {
	if updateTableCatelog == "all" {
		for catelog := range modules {
			err := services.UpdateMovies(catelog)
			if err != nil {
				return errors.New("更新" + catelog + "模块的信息失败: " + err.Error())
			}
		}
		fmt.Println("++ 更新所有模块成功!!!")
		fmt.Println("==========================")
	} else {
		err := services.UpdateMovies(updateTableCatelog)
		if err != nil {
			return errors.New("更新" + updateTableCatelog + "模块的信息失败: " + err.Error())
		}
	}
	return nil
}

//根据参数检查所有的文件,看哪些已经被下载，已经下载的视频会在数据库的记录中标记
func checkoutFiles(checkoutCatelog string) error {
	if checkoutCatelog == "all" {

	} else {
		err := services.CheckoutFiles()
	}
}
