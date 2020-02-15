###  uBeanos系统 VNC Server 配置

##### 1.1 系统源更新

```shell
sudo apt update
```

##### 1.2 安装`TigerVNC`

```shell
sudo apt install xserver-xorg-core

sudo apt install tigervnc-standalone-server tigervnc-xorg-extension tigervnc-viewer

sudo apt-get install gsfonts-x11 xfonts-base xfonts-75dpi xfonts-100dpi
```



##### 1.3 设置VNC密码

```shell
vncpasswd
```

输入密码，然后确认。



##### 1.4 启动脚本设置

我们将创建一个启动脚本作为初始配置，将在激活VNC服务器时执行

vim  ~/.vnc/xstartup

将下面代码拷入

```shell
#！/bin/sh

[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup

[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources

vncconfig -iconic &

dbus-launch --exit-with-session gnome-session &
```

保存退出，给~/.vnc/xstartup 文件添加可执行权限

```shell
sudo chmod a+x ~/.vnc/xstartup
```





##### 2.1 启动VN服务器

```shell
vncserver -localhost no -geometry 1920x950 -depth 24
```

上述选项将创建一个会话，允许外部连接具有1920x950像素分辨率和清晰度24.

大家可以更具自己喜好设置分辨率和清晰度



##### 2.2 连接VNC

大家可以使用VNC Viewer 连接到VNC服务器。

下载地址：  https://www.realvnc.com/en/connect/download/vnc/ 

vnc viewer安装好之后，点击左上角file菜单，选择new connection

注意： VNC server请填写主机ip地址+5901端口

例： 192.168.1.11::5901

name 和label 可以自定义填写





##### 3.1 管理VNC 服务器可能会使用的命令

查看vnc服务器的用户列表：

vncserver -list



查看VNC服务器是否处于活动状态：

pgrep Xtigervnc 或者 ss -tulpn | egrep -i 'vnc|590'



终止VNC会话

vncserver -kill  :1

上面的命令是终止ID为1的VNC会话，大家可以选择自己需要终止的会话。

