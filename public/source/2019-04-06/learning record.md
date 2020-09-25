title: 博客常用小功能
date: 2019-04-06 17:36:51
author: me
draft: false
cover: images/lr.png
top: false
type: post
hide: true
tags: #Optional
    - JavaScript
    - 网站
preview: Valine是一款优秀的评论系统，漂亮的界面、显示浏览器信息/系统信息、支持匿名评论、支持Markdown、Emoji等都是它的绝对优势.
-------
# 常用小功能实现(不定期更新......)
## 1、Valine  一款基于 Leancloud 的极简风评论系统
嗯，我想我需要评论功能，本博客设计方案是静态博客，不具有实时前后端数据交互功能(其实，压根就没有数据库，怎么实现动态数据交互......)，只能依赖于第三方评论系统。

之前有很多评论系统，包括Disqus，多说，畅言，网易云跟帖，Coding Comments…，但是这些评论系统，要不就是稳定性差，要不就是需要繁琐的登录，还有的直接停止服务了。

最初给博客配置的是畅言，奈何界面太丑，不定期弹广告，还需要备案，严重影响体验，果断放弃。


最终选择了Valine,Valine常用于Hexo,Hexo的设计与本博客不太一样，自己稍微改动了一下。
### 使用方法
#### 获取 Leancloud Key
因为是基于 Leancloud 的评论系统，所以需要先注册一个 Leancloud 账号。

- [点击这里注册 Leancloud 账号](https://leancloud.cn/dashboard/login.html#/signup)。
- 创建一个应用，应用名随意。
- 进入刚才创建的应用。设置 -> 应用 Key，这里的 App ID 和 App Key，复制下来，下一步要用。


#### 配置到页面上

### 创建文章
在`source`目录中建立任意`.md`文件（可置于子文件夹），使用如下格式：

``` 
<body>
<div class="comments" id="comments"></div>
    <script src="//cdn1.lncld.net/static/js/3.0.4/av-min.js"></script>
    <script src="//cdn.jsdelivr.net/npm/valine@1.1.6/dist/Valine.min.js"></script>
    <script>
        new Valine({
            av: AV,
            el: '.comments',
            notify: true, 
            verify: true,
            app_id: '复制的 App ID 粘贴到此处',
            app_key: '复制的 App Key 粘贴到此处',
            placeholder: '留下你的想法.....'
        });
    </script>
</body>
``` 


> 注：这种办法在我的博客中可以使用，其他博客系统不确定能否使用。
Hexo博客系统的配置方法自行百度。