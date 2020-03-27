title: 拥抱Linux
date: 2019-07-14 18:26:25
author: me
cover: /source/2019-05-03/timg.jpeg
draft: false
top: false
type: post
hide: false
tags: 
        - Linux
        - 记录
preview: 一直以来桌面使用的都是`Windows`，对`Windows`没有任何偏见，挺好的，没有网上流传那样不稳定，对`windows10`升级系统也不反感。也曾尝试过`VMware`里装`Ubuntn`, `Centos `，玩过树莓派(`Raspberry Pi`),就是玩玩没当回事，觉得`linux`还是属于少数人，理我还很远，`windows10`对我而言就足够了。但我在windows环境下把博客源码部署到服务的时候，我决定投入`linux`的怀抱。
---------------------
## 从Windows到Linux

Windows10用着挺好的，没有像网上流传的那样不稳定，不兼容，也没有遇到过蓝屏，Bug，强制更新之类的，用着还是蛮舒服的。只是目前`Windows`在某些方面不能满足我的需求，呃...也不是不能满足，是相比于Linux，Windows要麻烦一点，复杂一点，Linux一句命令能搞定的，Windows下要弄好久。

目前的情况而言，我还不能完全放弃Windows,用Linux来代替，毕竟在Linux下没有QQ，Office等常用软件，就算有替代品，体验也不是很好。(其实我是舍不得我的正版Windows和Office)，所以两全其美的办法就是装双系统！！！

## 安装Ubuntn	
### 准备工作
我是通过 U 盘制作启动盘安装的，所以整个过程只有一个前提，就是要准备一个容量大于 8G 的 U 盘，记得备份 U 盘内的重要文件，因为等会需要被格式化。

想要两个系统相互独立，则需要从 C 盘或者其它有空余空间的磁盘分割出一部分空间，我是 128G + 1T 双硬盘，我在机械硬盘分了 128G 出来，用于开发应该是够了。这个新分出来的磁盘在安装的时候是需要被设置为挂载点的。切记不要给它分配新盘符！！！


#### 一、准备工作要点

1、【Windows 10 上安装的软件】下载安装软碟通UltraISO，用于制作安装Ubuntu的U盘启动盘；

[UltraISO] (https://cn.ultraiso.net/) 使用方法很简单，自行百度。

2、 【Windows 10压缩Ubuntu所需空间】

在右键“此电脑”选择“管理”后的“磁盘管理”中，选择要压缩的卷，如　E 盘　“右键”－“压缩卷”－“102400” (压缩出来100Ｇ),此时，切记不要新建分区，就让它呈灰色状态，不管它。

#### 二、安装过程要点

1、【避免进选地图处安装程序卡死、安装失败的关键点】断网！！！

2、重新启动电脑，插入之前做的Ubuntu启动盘，从USB启动，进入Ubuntu启动界面

第一项是试用，点击后进入Ubuntu桌面系统，能够进行操作，但设置不能保存，可用于调试Linux系统，桌面的图标就是安装Ubuntu系统

第二项是直接安装Ubuntu系统，两种安装方式完全相同。

3、选择时间，语言，键盘，地区的自行选择以及你的姓名、计算机名、用户名和密码。

4、【安装类型】(重点)“安装类型”中选择“安装Ubuntu，与Windows10共存”。

5、点击现在重启，你将看到双系统选择界面。这时只能通过ubuntu进入windows了。

> - Grub会接管启动，并列出可以启动项目，通常Ubuntu在第一项，Ubuntu高级设置在第二项，Windows Boot Manager在第三项，您有10秒的选择时间。

## 后续

#### 1 、WIindows 10 与Ubuntn 时间不同步问题
Ubuntu系统下，打开终端(Ctrl + Alt + t)输入：
```
sudo timedatectl set-local-rtc 1
```
下面更新一下时间
```
$ sudo apt-get install ntpdate
$ sudo ntpdate time.windows.com
```
将时间更新至硬件
```
$ sudo hwclock --localtime --systohc
```
重启进入Windows，发现Windows时间正常

#### 2、修改默认启动，以默认启动　Windows 10 为例
Ubuntu系统下，打开终端(Ctrl + Alt + t)输入：
```
sudo gedit /etc/default/grub
```
输入你的密码以后会自动打开文件


``` yml 
# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.
# For full documentation of the options in this file, see:
#   info -f grub -n 'Simple configuration'

GRUB_DEFAULT=2    /* 修改这个参数　*/
GRUB_TIMEOUT_STYLE=hidden
GRUB_TIMEOUT=10
GRUB_DISTRIBUTOR=`lsb_release -i -s 2> /dev/null || echo Debian`
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_CMDLINE_LINUX=""

# Uncomment to enable BadRAM filtering, modify to suit your needs
# This works with Linux (no patch required) and with any kernel that obtains
# the memory map information from GRUB (GNU Mach, kernel of FreeBSD ...)
#GRUB_BADRAM="0x01234567,0xfefefefe,0x89abcdef,0xefefefef"

# Uncomment to disable graphical terminal (grub-pc only)
#GRUB_TERMINAL=console

# The resolution used on graphical terminal
# note that you can use only modes which your graphic card supports via VBE
# you can see them in real GRUB with the command `vbeinfo'
#GRUB_GFXMODE=640x480

# Uncomment if you don't want GRUB to pass "root=UUID=xxx" parameter to Linux
#GRUB_DISABLE_LINUX_UUID=true

# Uncomment to disable generation of recovery mode menu entries
#GRUB_DISABLE_RECOVERY="true"

# Uncomment to get a beep at grub start
#GRUB_INIT_TUNE="480 440 1"
```
> 其中GRUB_DEFAULT=0代表系统默认启动第0项，因为我的windows启动项是第三项，又因为是从0开始计算(0 , 1 , 2 , 3)，所以将它改成2， 
GRUB_TIMEOUT=10，代表的是选择时间，这里我将不修改，保存并退出

```
sudo update-grub
```

此时重启电脑，会看到选择条自动选中　Windows Boot Manager (on /dev/sdb1)，可是用键盘上下键选择要启动的系统，若不做调整，10秒后默认启动　Windows 10 。
	