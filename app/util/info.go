package util

import (
	"fmt"
	"os"
)

const (
	CLR_R = "\x1b[31;1m"
	CLR_Y = "\x1b[33;1m"
)

//自定信息输出
func MLog(info interface{}) {
	fmt.Printf("%s\n", info)
}

//打印警告不退出
func MWarn(info interface{}) {
	fmt.Printf("WARNING: %s%s\n%s", CLR_Y, info, "\x1b[0m")
}

//只打印错误不退出
func MError(info interface{}) {
	fmt.Printf("ERR: %s%s\n%s", CLR_R, info, "\x1b[0m")
}

//打印错误日志并退出
func MFatal(info interface{}) {
	MError(info)
	os.Exit(1)
}
