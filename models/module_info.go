package models

import (
	"cn/com/hengwei/commons/cfg"
	"path/filepath"
)

type ModuleInfo struct {
	Catelog      string
	Name         string
	FilePath     string
	MovieDirPath string
}

//要爬取的小电影主页
var HttpPath string

//存放小电影的文件夹路径(大路径)
var DirPath string

//初始化
func InitModuleInfos() (map[string]ModuleInfo, error) {
	infos := []ModuleInfo{
		ModuleInfo{"60", "国产精品", "", ""},
		ModuleInfo{"62", "欧美性爱", "", ""},
		ModuleInfo{"87", "夫妻同房", "", ""},
		ModuleInfo{"88", "手机小视频", "", ""},
		ModuleInfo{"89", "自拍偷拍", "", ""},
		ModuleInfo{"89", "换妻游戏", "", ""},
		ModuleInfo{"91", "网红主播", "", ""},
		ModuleInfo{"92", "明星艳照门", "", ""},
		ModuleInfo{"93", "开放90后", "", ""},
		ModuleInfo{"101", "成人动漫", "", ""},
		ModuleInfo{"110", "亚洲无码", "", ""},
		ModuleInfo{"130", "高清无码", "", ""},
	}

	modules := map[string]ModuleInfo{}
	for _, info := range infos {
		info.FilePath = filepath.Join(DirPath, info.Name+".txt")
		info.MovieDirPath = filepath.Join(DirPath, info.Name)
		modules[info.Catelog] = info
	}

	return modules, nil
}

func init() {
	paths, err := cfg.ReadProperties(filepath.Join("conf", "path.properties"))
	if err != nil {
		panic("读取path.properties失败： " + err.Error())
	}
	DirPath = paths["fileDirPath"]
	HttpPath = paths["httpPath"]
}
