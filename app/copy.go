package app

import (
	"io"
	"os"
	"path/filepath"
)

//总的复制入口
func Copy() {
	srcList := globalConfig.Build.Copy
	for _, source := range srcList {
		//多层目录查询
		if matches, err := filepath.Glob(filepath.Join(rootPath, source)); err == nil {
			for _, srcPath := range matches {
				MLog("Copying: " + srcPath)
				file, err := os.Stat(srcPath)
				if err != nil {
					MLog("没有找到文件：" + srcPath)
				}
				fileName := file.Name()

				//目标粘贴目录
				desPath := filepath.Join(publicPath, fileName)
				WaitGroup.Add(1)
				// 文件，或目录的复制
				if file.IsDir() {
					go CopyDir(srcPath, desPath)
				} else {
					go CopyFile(srcPath, desPath)
				}
			}
		} else {
			MFatal(err.Error())
		}
	}
}

//复制文件夹
func CopyDir(source string, dest string) {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		MFatal(err.Error())
	}
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		MFatal(err.Error())
	}
	directory, _ := os.Open(source)
	defer directory.Close()
	defer WaitGroup.Done()
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		//源文件指针
		sourcefilepointer := source + "/" + obj.Name()
		//目标文件指针
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			WaitGroup.Add(1)
			CopyDir(sourcefilepointer, destinationfilepointer)
		} else {
			WaitGroup.Add(1)
			go CopyFile(sourcefilepointer, destinationfilepointer)
		}
	}
}

//复制文件
func CopyFile(source string, dest string) {
	sourcefile, err := os.Open(source)
	defer sourcefile.Close()
	if err != nil {
		MFatal(err.Error())
	}
	destfile, err := os.Create(dest)
	if err != nil {
		MFatal(err.Error())
	}
	defer destfile.Close()
	defer WaitGroup.Done()
	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
}
