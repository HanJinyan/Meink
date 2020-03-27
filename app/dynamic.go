package app

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

/*
动态监测文件，目录变化 ，实时编译
*/
func DynamicMonitoringFile() {
	var dynamic *fsnotify.Watcher
	if dynamic != nil {
		dynamic.Close()
	}
	dynamic, _ = fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case event := <-dynamic.Events:
				if event.Op == fsnotify.Write {
					MLog(event.Name)
					ParseGlobalConfigForWrap(true)
					Build()
				}
			case err := <-dynamic.Errors:
				MFatal(err.Error)
			}
		}
	}()
	var dirs = []string{
		filepath.Join(rootPath, "source"),
		filepath.Join(themePath, "bundle"),
	}
	var files = []string{
		filepath.Join(rootPath, "config.yml"),
		filepath.Join(themePath),
	}

	for _, source := range dirs {
		SymWalk(source, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				if err := dynamic.Add(path); err != nil {
					MFatal(err.Error)
				}
			}

			return nil
		})
	}

	for _, source := range files {
		if err := dynamic.Add(source); err != nil {
			MFatal(err.Error())
		}
	}
}
