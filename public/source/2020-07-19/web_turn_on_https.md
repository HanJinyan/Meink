title: 静态博客开启HTTPS
date: 2020-07-19 10:20:57
author: me
cover: /source/2020-07-19/https.png
draft: false
top: false
type: post
hide: false
tags:
    -  HTTPS
    -  网站

preview: HTTP报文以明文形式传输，存在一定的安全隐患，有可能遭受到安全攻击。使用 `Google`浏览器打开一个 `HTTP`协议网站，会发现 `Chrome`在网址的左边将这个网站标记为不安全。嗯.......总感觉缺点啥，不爽，怎么能容许自己与世界脱轨，我也搞一个吧。
---------------------

### 前言

#### 什么是HTTPS

> 百度百科：
>
> HTTPS （全称：Hyper Text Transfer Protocol over SecureSocket Layer），是以安全为目标的 HTTP 通道，在HTTP的基础上通过传输加密和身份认证保证了传输过程的安全性 。HTTPS 在HTTP 的基础下加入SSL层，HTTPS 的安全基础是 SSL，因此加密的详细内容就需要 SSL。 HTTPS 存在不同于 HTTP 的默认端口及一个加密/身份验证层（在 HTTP与 TCP之间）。这个系统提供了身份验证与加密通讯方法。它被广泛用于万维网上安全敏感的通讯，例如交易支付等方面

#### 说人话

HTTPS 为 HTTP 报文提供了一个加密传输的通道，这样攻击者就无法窃听或者篡改传输的内容。要启用 HTTPS，必须向一个可信任机构申请一个 HTTPS 证书。专业的证书申请需要收费，不过对于个人博客网站来说，有很多免费的证书申请机构，当然，家里有矿的读者嘛，你可以关闭这篇文章了，免费的证书申请机构比如Let’s Encrypt，它提供了免费的证书申请服务，申请过程十分简单，只需要运行几条命令即可，而且证书到期后支持自动续期，可谓一劳永逸。接下来我们就是用 Let’s Encrypt 提供的工具来申请免费的 HTTPS 证书。

#### 一、准备工作

如果你的网站用的是80端口通信，想必你用的是HTTP，你需要`修改`你的对外通信`端口`，这个是必须的！！！因为HTTPS用的是`443`端口。以小编的博客为例：

- 1、登陆阿里云

- 2、选择控制台 -- 云服务器 ECS -- 实例与镜像 -- 实例

- 3、选择你的实例进入 -- 本实例安全组 -- 配置规则

- 4、入方向 -- 手动添加

- 5、端口范围：443，源：0.0.0.0/0 ，其他默认

- 6、保存退出

> 注：想要用HTTPS，你需要使用Nginx ，建议做端口转发。

#### 二、服务器端操作

登录[Certbot](https://certbot.eff.org/)选择我们博客网站使用的服务器软件和操作系统，小编以 Nginx 和 Ubuntn 为例：

##### 1、ssh 登陆你的服务器

##### 2、添加Certbot PPA

```shell
sudo apt-get update
sudo apt-get install software-properties-common
sudo add-apt-repository universe
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
```

> 上面命令依次执行

##### 3、安装Certbot

```shell
sudo apt-get install certbot python3-certbot-nginx
```

##### 4、获取并安装证书

```shell
sudo certbot --nginx
```

> 注：会有一系列交互式的提示，首先会让你输入**邮箱**，用于订阅。然后输入**a** 同意他们的政策。

接着 certbot 会自动扫描出来域名，根据提示输入想开启 HTTPS 的域名的**标号**，

例：

````
Which names would you like to activate HTTPS for

1:www.jinyan.ink

2:www.meink.jinyan.ink

Select the appropriate numbers separated by commas and/or spaces, 
or leave input blank to select all options shown (Enter 'c' to cancel): 1
````

然后 certbot 会做一个域名校验，证明你对这个域名有控制权限。验证通过后，Let's Encrypt 就会把证书颁发给你。最后会提示你是否把 HTTP 重定向到 HTTPS，选择**是**（数字标号2），这样 certbot 会自动帮我们修改 Nginx 的配置，将 HTTP 重定向到 HTTPS，如果用户使用 HTTP 协议访问我们的博客网站，就会重定向到 HTTPS 协议访问，确保安全性。

例：

```
Please choose whether or not to redirect HTTP traffic to HTTPS, removing HTTP access.

1: No redirect - Make no further changes to the webserver configuration.

2: Redirect - Make all requests redirect to secure HTTPS access. Choose this for

new sites, or if you're confident your site works on HTTPS. You can undo this

change by editing your web server's configuration.

Select the appropriate number [1-2] then [enter] (press 'c' to cancel): 2 
```

##### 5、自动自动续期

certbot 申请的证书只有 3 个月有效期，不过没有关系，certbot 可以无限续期，我们增加一条 crontab 定时任务用来执行 certbot 自动续期任务，这样一次申请，终生使用。

打开 /etc/crontab，增加定时任务:

```
echo "0 0,12 * * * root python -c 'import random; import time; time.sleep(random.random() * 3600)' && certbot renew" | sudo tee -a /etc/crontab > /dev/null
```

> 这里配置每天 12 点执行自动续期命令。

### 后续

##### 1、 全局HTTPS

由于全站开启了 HTTPS，因此需要把网站中非 HTTPS 的内容（比如通过 HTTP 协议请求的外部资源）改为 HTTPS，改为https://www.example.com

##### 2、Websock 通信

需要把：

```html
WebSocket('ws://' + location.host);
```

修改为：

```
WebSocket('wss://' + location.host);
```