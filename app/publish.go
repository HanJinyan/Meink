package app

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Release() {
	RunCommand(globalConfig.Publish.Release)
	MLog(globalConfig.Publish.Release)
}
func Sync(){
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
					RunCommand(globalConfig.Publish.Sync)
				}
			case err := <-dynamic.Errors:
				MFatal(err.Error)
			}
		}
	}()
	var dirs = []string{
		filepath.Join(rootPath, "source"),
	}
	var files = []string{
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


func RunCommand(command string) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	cmd := exec.Command(shell,flag,command)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %s......", err.Error())
		return err
	}
	go rsyncLog(stdout)
	go rsyncLog(stderr)
	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command execution: %s......", err.Error())
		return err
	}
	return nil
}
func rsyncLog(reader io.ReadCloser) error {
	cache := "" //缓存不足一行的日志信息
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err!=io.EOF{
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
			fmt.Printf("%s%s\n", cache, line)
			cache = s[len(s)-1]
		}
	}
	return nil
}