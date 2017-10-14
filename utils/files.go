package utils

import (
	"io/ioutil"
	"path/filepath"
)

func GetFilesNameByDirPath(dirPath string) []string {
	names := []string{}
	getFilesNameByDirPath(dirPath, &names)
	return names
}

//两个参数 分别表示：文件夹路径 结果只包含文件  不包含文件夹（包含子文件夹里面的文件）
func getFilesNameByDirPath(dirPath string, names *[]string) {
	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic("读取" + dirPath + "失败: " + err.Error())
	}
	for _, info := range infos {
		if info.IsDir() {
			newPath := filepath.Join(dirPath, info.Name())
			getFilesNameByDirPath(newPath, names)
		} else {
			*names = append(*names, info.Name())
		}
	}
}
