title: Golang 执行 shell脚本并实时打印 shell 脚本输出日志信息
date: 2020-04-09 15:29:08
author: me
cover: 
draft: false
top: false
type: post
hide: false
tags: 
    - Golang
preview: Meink博客自动同步文章到服务器要用到 shell 脚本的配合来完成，需要实时打印 shell 脚本的中每条命令的输出日志信息，便于查看任务进度等,但 `Golang` 自带的 `cmd`库可以完成命令的执行，但不能够打印执行过程 。
---------------------
## Golang 执行 shell脚本，并实时打印 shell 脚本输出日志信息
博客自动同步文章到服务器要用到 shell 脚本的配合来完成，需要实时打印 shell 脚本的中每条命令的输出日志信息，便于查看任务进度等
但 `Golang` 自带的 `cmd`库可以完成命令的执行，但不能够打印执行过程。

### 源码如下


````
package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"fmt"
)

func main() {
	RunCommand("ls -lh")   //替换为你的脚本
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
````