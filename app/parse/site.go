package parse

import (
	"Meink/app/modle"
	"Meink/app/system"
	"Meink/app/util"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var siteConfig *modle.SiteConfig

//SiteConfig 解析 config.yml 配置信息
func SiteConfig() *modle.SiteConfig {
	runPath := system.RunPath()
	appConfig := AppConfig()
	data, err := ioutil.ReadFile(filepath.Join(runPath, "config.yml"))

	if err != nil {
		util.MFatal("没有找到 config.yml 这个文件 ，或者文件为空: " + err.Error())
	}

	if err = yaml.Unmarshal(data, &siteConfig); err != nil {
		util.MFatal("解码config.yml 文件出错:" + err.Error())
	}
	siteConfig.Site.Logo = strings.Replace(siteConfig.Site.Logo, "-/", siteConfig.Site.Logo, 1)

	if siteConfig.Site.URL != "" && strings.HasSuffix(siteConfig.Site.URL, "/") {
		siteConfig.Site.URL = strings.TrimSuffix(siteConfig.Site.URL, "/")
	}

	if siteConfig.Build.Output == "" {
		siteConfig.Build.Output = "public"
	}

	if siteConfig.Build.Source == "" {
		siteConfig.Build.Source = "source"
	}
	themeConfig := ThemeConfig(filepath.Join(runPath, siteConfig.Site.Theme, "config.yml"))
	for _, copyItem := range themeConfig.Copy {
		siteConfig.Build.Copy = append(siteConfig.Build.Copy, filepath.Join(siteConfig.Site.Theme, copyItem))
	}
	siteConfig.App.Version = appConfig.Version
	siteConfig.App.Git = appConfig.Git
	siteConfig.App.Author = appConfig.Author
	siteConfig.App.Email = appConfig.Email
	siteConfig.App.Name = appConfig.Name
	siteConfig.App.Port = appConfig.Port
	return siteConfig
}
