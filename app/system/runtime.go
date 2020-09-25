package system

import (
	"Meink/app/util"
	"os"
)

func RunPath() string {
	dir, err := os.Getwd()
	if err != nil {
		util.MFatal("获取当前可执行路径错误:" + err.Error())
	}
	return dir
}
