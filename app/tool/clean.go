package tool

import (
	"Meink/app/parse"
	"Meink/app/system"
	"Meink/app/util"
	"os"
	"path/filepath"
)

//清空public文件夹
//Debug的时候方便查看文件是否生成
func CleanPublic() {
	runPath := system.RunPath()
	cleanPath := parse.SiteConfig().Build.Output
	publicPath := filepath.Join(runPath, cleanPath)
	//要清理的文件夹 ，后缀名
	cleanPatterns := []string{"tag", "bundle", "misc", "images", "source", "*.html", "*.ico", "*.png", "*.txt", "*.xml"}
	for _, pattern := range cleanPatterns {
		files, _ := filepath.Glob(filepath.Join(publicPath, pattern))
		for _, path := range files {
			os.RemoveAll(path)
			util.MLog("Cleaning " + path)
		}

	}
}
