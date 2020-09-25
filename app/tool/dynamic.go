package tool

import (
	"Meink/app/build"
	"Meink/app/parse"
	"Meink/app/system"
	"Meink/app/util"

	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

/*
动态监测文件，目录变化 ，实时编译
*/
func DynamicMonitoringFile() {
	var dynamic *fsnotify.Watcher
	runPath := system.RunPath()
	siteConfig := parse.SiteConfig()
	themePath := filepath.Join(runPath, siteConfig.Site.Theme)

	if dynamic != nil {
		dynamic.Close()
	}
	dynamic, _ = fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case event := <-dynamic.Events:
				if event.Op == fsnotify.Write {
					util.MLog(event.Name)
					build.Build()
				}
			case err := <-dynamic.Errors:
				util.MFatal(err.Error)
			}
		}
	}()
	var dirs = []string{
		filepath.Join(runPath, siteConfig.Build.Source),
		filepath.Join(themePath, "bundle"),
	}
	var files = []string{
		filepath.Join(runPath, "config.yml"),
		filepath.Join(themePath),
	}

	for _, source := range dirs {
		util.SymWalk(source, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				if err := dynamic.Add(path); err != nil {
					util.MFatal(err.Error)
				}
			}

			return nil
		})
	}

	for _, source := range files {
		if err := dynamic.Add(source); err != nil {
			util.MFatal(err.Error())
		}
	}
}
