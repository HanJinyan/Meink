title: Golang 跨平台批量编译脚本 
date: 2020-05-04 20:48:00
author: me
cover: 
draft: false
top: false
type: post
hide: false
tags: 
    - Golang
    - Shell
preview: 解决项目打包发布时，在本地批量编译多个平台的可执行文件。
---------------------
### Golang 跨平台批量编译脚本 
````
#!/bin/bash
build() {
  echo "Building for $1 $2..."
  suffix=""
  if [ $1 == "windows" ]; then
    suffix=".exe"
  fi
  GOOS=$1 GOARCH=$2 go build
}
build windows 386
build windows amd64
build windows arm
build linux 386
build linux amd64
build linux arm
build linux arm64
build darwin 386
build darwin amd64
````