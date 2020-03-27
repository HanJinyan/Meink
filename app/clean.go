package app

import (
	"os"
	"path/filepath"
)

//清空public文件夹
//Debug的时候方便查看文件是否生成
func CleanPublic() {
	publicPath = filepath.Join(rootPath, globalConfig.Build.Output)
	//要清理的文件夹 ，后缀名
	cleanPatterns := []string{"tag", "bundle", "misc", "images", "*.html", "*.ico", "*.png", "*.txt"}
	for _, pattern := range cleanPatterns {
		files, _ := filepath.Glob(filepath.Join(publicPath, pattern))
		for _, path := range files {
			os.RemoveAll(path)
			MLog("Cleaning " + path)
		}

	}
}
