title: Linux Command
date: 2020-03-17 20:14:10
author: me
cover: 
draft: false
top: false
type: post
hide: false
tags:   
    -  命令 
preview: 常用的一些命令，总是记不住，写这里把。
---------------------
## Linux Command 
#### 查询应用占用的端口
```
netstat -nap | grep 8080
```
```
netstat -nap | grep Meink
```
#### Ubuntn 修改更新版本 如 lts ,normal ,naver
```
sudo gedit /etc/update-manager/release-upgrades
```
#### 检查Ubuntn 是否有新的lts 版本
```
sudo do-release-upgrade
```
#### 检查Ubuntn 调用更新器更新
```
sudo update-manager -c -d
```
#### 初次git上传
```
ssh-keygen -t rsa -C "your email"

cat ~/.ssh/id_rsa.pub

git config --global user.name "your name"

git config --global user.email  "your email"

git add .

git commit -m "update"

git push -u origin master 

git push -f origin master  //强制上传 ，有风险

```
cat /proc/pid/status

