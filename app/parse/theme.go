package parse

import (
	"Meink/app/modle"
	"Meink/app/system"
	"Meink/app/util"
	"io/ioutil"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

//ThemeConfig 解析主题语言配置，复制目录
func ThemeConfig(fileName string) *modle.Theme {
	var themeConfig *modle.Theme
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		util.MFatal("theme里面没有找到config.yml文件:" + err.Error())
	}
	if err := yaml.Unmarshal(data, &themeConfig); err != nil {
		util.MFatal("theme里面的onfig.yml文件存在错误:" + err.Error())
	}
	return themeConfig
}

//这个函数极度不安全，并发情况下偶尔是出错，没有找到好的解决方案
func I18n() map[string]string {
	var lock sync.Mutex
	i18n := make(map[string]string)
	runPath := system.RunPath()
	themeConfig := ThemeConfig(filepath.Join(runPath, siteConfig.Site.Theme, "config.yml"))

	for item, languageItem := range themeConfig.Language {
		lock.Lock()
		i18n[item] = languageItem[siteConfig.Site.Language] //languageItem map[string]string
		lock.Unlock()
	}

	return i18n
}

// func I16N() sync.Map {
// 	var i16n sync.Map
// 	runPath := system.RunPath()
// 	themeConfig := ThemeConfig(filepath.Join(runPath, siteConfig.Site.Theme, "config.yml"))
// 	for item, languageItem := range themeConfig.Language {
// 		i16n.Store(item, languageItem[siteConfig.Site.Language])
// 	}
// 	return i16n

// }
// func RR() {
// 	ss := I16N()
// 	f := func(k, v interface{}) bool {
// 		return true
// 	}
// 	ss.Range(f)
// }
