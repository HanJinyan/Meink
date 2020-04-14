#!/bin/bash
# 需要自己修改部分文件目录
build() {
  echo "Building for $1 $2..."
  rsync -av ./theme source config.yml release/Meink_$1_$2
  suffix=""
  if [ $1 == "windows" ]; then
    suffix=".exe"
  fi
  GOOS=$1 GOARCH=$2 go build -ldflags="-w -s" -o release/Meink_$1_$2/Meink_$1_$2$suffix
  cd release
  tar cvf - Meink_$1_$2/* | gzip -9 - >Meink_$1_$2.tar.gz
  rm -rf Meink_$1_$2
  cd ..
}
rm -rf release
mkdir -p release
build windows 386
build windows amd64
build windows arm
build linux 386
build linux amd64
build linux arm
build linux arm64
build darwin 386
build darwin amd64