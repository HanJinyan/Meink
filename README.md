### Meink是一个基于Golang的轻量级静态博客构建工具，可以快速搭建博客网站。它无依赖跨平台，配置简单构建快速，以`Markdown`格式文本输入源，注重简洁易用.
快速预览，可以关注我的个人博客：www.jinyan.ink 项目源码整体比较简单，适合初学者练手学习项目。

### 特点
- YAML格式的配置
- Markdown格式的文章
- 无依赖跨平台
- 超快的构建速度(并发)
- 简单，小巧，易用
- 源码简单，适合初学者练手

## 快速开始
### 下载，安装Golang
``` 
go version go1.14 linux/amd64
```

>- 注：go mod 模式开启

### 下载源码
``` 
go get github.com/HanJinyan/Meink
```
### 编译源码并运行 (Linux)
```
hanjinyan@hanjinyan:~/app/src/Meink$ go build
hanjinyan@hanjinyan:~/app/src/Meink$ ./Meink run
```
浏览器打开：localhost:8080　　进行预览

### 命令
    $ ./Meink run      运行博客
    $ ./Meink init     清空public (用于调试)
    $ ./Meink new      创建文章

### 网站配置

编辑config.yml，使用格式：

```
site:
  title: "Meink· 轻量级静态博客构建工具"   # 网站名称
  subtitle: "Keep It Simple & Stupid"   # 网站首页字幕
  limit: 5                              # 每页文章数
  theme: theme                          # 主题名称（默认是theme）
  language: zh-cn                       # 网站显示语言(可选zh-cn ,en )
  url: "http://www.jinyan.ink"          # 网站url
  root: "/Meink"                        # 网站根目录
  copyright: "你的备案号"                # 网站备案号
  email: "你的邮箱" 
  github: "你的Github名称"
  link: "{year}{month}{day}{hour}{minute}{second}.html" #文章连接


authors:
  me:
    name: "Meink"                       # 文章作者名
    intro: "Meink"                      # 文章作者描述
    avatar: "/author.png"               # 文章作者logo


build:
    output: "public"                    # 输出目录
    source: "source"                    # 文章存放目录
    port: 8080                          # 输出端口
    copy:
        - "source"
```
### 写作
#### 用命令创建

```
./Meink new blog_name  # 创建文章
```
>-  blog_name.md文件在source/xxxx年-xx月-xx日 文件夹下

#### 文章格式 

编辑.md文件，使用格式： 

```
title: Meink                    # 文章名称
date: 2020-03-27 18:47:26       # 文章时间(自动生成)
update: 2020-03-27 18:47:26     # 文章更新时间(手动生成)
author: me                      # 文章作者(默认me)
cover:                          # 文章封面
draft: false                    # 是否为草稿(草稿不渲染)
top: false                      # 是否置顶
type: post                      # 文章类型(post为文章 ， page 为页面)
hide: false                     # 是否隐藏(隐藏文章不渲染)
tags:                           # 文章标签
    - Meink
    - GitHub
preview:                        # 文章预览内容
---------------------
Markdown 格式的正文
```
### 用命令创建页面

```
./Meink new page_name  # 创建页面
```
#### 页面格式 

编辑.md文件，使用格式： 

```
type: page                      # 类型(post为文章 ， page 为页面)
title: "关于"                    # 页面标题       
---
Markdown 格式的正文
```

#### 手动创建

>-  .md创建在source文件夹下 (允许子文件)

### 二次开发

 欢迎各位爱好者提建议，二次开发，源码交流

### Bug 
  解析语言过程中使用了Map结构，在渲染并发过程中有几率触发 "concurrent map writes and read" 这个Bug 导致程序挂掉，目前自己没能解决。

### 联系我

> QQ : 994205825

> 微信： HBY205825

> 邮箱：hby0210@163.com