package build

import (
	"Meink/app/util"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func MergeTpl(htmlPath string) string {
	var partialTpl string
	files, _ := filepath.Glob(filepath.Join(htmlPath, "*.html"))
	for _, path := range files {
		fileExt := strings.ToLower(filepath.Ext(path))
		baseName := strings.ToLower(filepath.Base(path))
		if fileExt == ".html" && strings.HasPrefix(baseName, "_") {
			html, err := ioutil.ReadFile(path)
			if err != nil {
				util.MFatal(err.Error())

			}
			tplName := strings.TrimPrefix(baseName, "_")
			tplName = strings.TrimSuffix(tplName, ".html")
			htmlStr := "{{define \"" + tplName + "\"}}" + string(html) + "{{end}}"
			partialTpl = partialTpl + htmlStr

		}

	}
	return partialTpl
}
