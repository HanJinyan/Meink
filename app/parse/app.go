package parse

import (
	"Meink/app/modle"
	"Meink/app/system"
	"Meink/app/util"

	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//AppConfig 输出解析APP 信息
func AppConfig() *modle.App {
	runPath := system.RunPath()
	var appConfig modle.App
	data, err := ioutil.ReadFile(filepath.Join(runPath, "/app/conf/config.yml"))
	if err != nil {
		util.MFatal("没有找到 config.yml 这个文件 ，或者文件为空: " + err.Error())
	}
	if err = yaml.Unmarshal(data, &appConfig); err != nil {
		util.MFatal("解码config.yml 文件出错:" + err.Error())
	}
	return &appConfig
}
