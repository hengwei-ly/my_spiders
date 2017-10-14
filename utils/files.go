package utils

import (
	"errors"
	"io/ioutil"
	"os"
)

//获取文件夹下的所有文件信息
func GetFilesInfoByDirPath(dirPath string) ([]os.FileInfo, error) {
	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, errors.New("")
	}
	return infos, err
}
