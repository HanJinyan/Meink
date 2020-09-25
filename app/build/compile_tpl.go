package build

import (
	"Meink/app/parse"
	"Meink/app/util"
	"html/template"
	"io/ioutil"
)

//连接并编译模板
func CompileTpl(tplPath string, partialTpl string, name string) template.Template {
	htmlTpl, err := ioutil.ReadFile(tplPath)
	if err != nil {
		util.MFatal("连接并编译模板出错" + err.Error())
	}
	//把模板合并
	htmlStr := string(htmlTpl) + partialTpl
	//插入 I18n数据
	funcMap := template.FuncMap{
		"i18n": func(val string) string {
			return parse.I18n()[val]
		},
	}
	tpl, err := template.New(name).Funcs(funcMap).Parse(htmlStr)
	if err != nil {
		util.MFatal(err.Error())
	}

	return *tpl
}
