title: Meink 设计构思(持续更新)
date: 2020-02-29 23:02:24
author: me
cover:
draft: false
top: false
type: post
hide: false
tags: 
    -  Meink
    -  源码
preview:  Meink 设计构思及源码解读
---------------------
### Meink 设计构思
#### 基本文件视图
``` 
.Meink
├── app            //Go实现的静态博客生成器
│   ├── build.go
│   ├── clean.go
│   ├── const.go
│   ├── container.go
│   ├── copy.go
│   ├── core.go
│   ├── dynamicfile.go
│   ├── info.go
│   ├── new.go
│   ├── parse.go
│   ├── publish.go
│   ├── rander.go
│   ├── serve.go
│   ├── symwalk.go
│   └── websocket.go
├── public         //前端文件渲染目录
├── source         //文章保存目录
│   ├── 2019-02-24
│   │   └── welcome.md
│   ├── about.me.md
│   └── images
│       └──hpfavicon.png
│  
├── theme          //主题目录 
│   ├── bundle     //静态文件目录 .css .js .img
│   │   ├── 20190302133942.jpg
│   │   ├── github-fill.svg
│   │   ├── home-bg.jpg
│   │   └── index.css
│   ├── _footer.html
│   ├── _header.html
│   ├── _head.html
│   ├── archive.html
│   ├── article.html
│   ├── page.html
│   ├── robots.txt
│   ├── tag.html
│   └── config.yml  //主题配置文件 .i18n的语言支持
│ 
│ 
├── config.yml      //站点配置
└── Meink.go        //Go入口
```

### 设计构思
#### 基本流程如下：
```

.md源文章 ---> Meink ---> .HTML文件

```
- 1、Meink 读取 .md 源文章

- 2、Meink 读取配置信息和 .md 文章内容

- 3、Meink 将读取的配置信息和文章内容并插入 .HTML 模板中

- 4、Meink 构建生成对应的 .HTML 文件

- 5、Meink 渲染 .HTML文件

>- 整个项目的重点是 `Meink` ，同时`Meink`也是核心。

### 源码
#### 环境

```
go version go1.14 linux/amd64
```
>- go mod 开启

#### 源码解析
- 对于 `Golang`，`Git`, 开发环境的安装配置不做过多的赘述。

- 同时只对核心部分源码解析。

#### 起步

参考基本文件视图，建好对应的文件，文件夹




