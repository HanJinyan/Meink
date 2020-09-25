package build

import (
	"Meink/app/util"
	"html/template"
	"os"
)

//按数据呈现html文件
func RenderPage(tpl template.Template, tplData interface{}, outPath string) {
	outFile, err := os.Create(outPath)
	if err != nil {
		util.MFatal(err.Error())
	}
	defer func() {
		outFile.Close()
	}()
	defer WaitGroup.Done()
	// 模板渲染
	err = tpl.Execute(outFile, tplData)
	if err != nil {
		util.MFatal(err.Error())
	}
}
