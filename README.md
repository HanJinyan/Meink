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
    $ ./Meink publish  部署到服务器
    $ ./Meink new      创建文章

### 二次开发

 欢迎各位爱好者提建议，二次开发，源码交流
 
### 联系我

> QQ : 994205825

> 微信： HBY205825

> 邮箱：hby0210@163.com