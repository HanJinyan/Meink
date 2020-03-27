#!/bin/sh
mkdir -p release
rsync -av ./* release/Meink release/Meink --delete --exclude public
GOOS=linux GOARCH=amd64 go build
cd release
rm -rf release
tar cvf - Meink/* | gzip -9 - > Meink.tar.gz
