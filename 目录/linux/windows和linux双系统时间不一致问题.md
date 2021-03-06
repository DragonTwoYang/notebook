#### Windows和Linux双系统时间不一致问题

------

##### 一. 时间的几个基本概念

首先我们需要先了解系统上时间几个基本时间概念

- GMT  即格林尼治标准时间，也就是世界时。GMT的正午是指当太阳横穿格林尼治子午线（本初子午线）时的时间。但由于地球自转不均匀不规则，导致GMT不精确，现在已经不再作为世界标准时间使用。

- UTC   即协调世界时。UTC是以原子时秒长为基础，在时刻上尽量接近于GMT的一种时间计量系统。为确保UTC与GMT相差不会超过0.9秒，在有需要的情况下会在UTC内加上正或负闰秒。UTC现在作为世界标准时间使用。
- LocalTime(本地时间，本地时间就是UTC+时区)

现在基本不使用GMT时间，但我们其实可以认为GMT和UTC时间是一个时间，相差无几。

------

##### 二. 双系统时间不一致的原因

由于电脑随时可能会处于断电和断网的情况，这个时候，操作时间保存时间就是一个问题。

由于上面的原因，前人在BIOS中保存了一份时间，但大家对于保存的时间是UTC时间和LocalTime时间没有做规定。这个由操作系统自行决定。

现在windows 默认BIOS中保存的是LocalTime时间。

linux 默认BIOS中保存的是UTC时间。

大家可以看到，这个就是我们安装双系统的时候，时间经常差8个小时的原因。因为大家默认使用北京时间，东八区。

基于上面所说，我们其实有2个办法来改正这个问题：

1. 将windows 默认BIOS时间由LocalTime 改为UTC

2. 将linux 默认BIOS时间由UTC改为 LocalTime

   *记住：上面2个方法只需要改一个即可。*

------

##### 三.  解决方案

先说windows吧，windows把配置都集成在注册表中，我们只需要修改注册表后重启电脑就可以搞定：

在 HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\TimeZoneInformation\ 中
添加一项数据类型为 REG_DWORD 并将值设为1即可

linux下面可以使用timedatectl 命令

```shell
timedatectl set-local-rtc 1
```

------

##### 四. 新的问题

这个问题一般大家不太可能碰到，我先描述下我碰到的问题：

我们公司有一款制作数据镜像的工具，使用ubuntu作为livecd，每次制作数据都从livecd中启动，但发现一个问题，每次运行完livecd之后，源目标的windows系统时间会出错，且linux下面时间也不对。

这个问题其实就是上面说的对硬件时间认识不一致导致的。

网上搜索了下，大部分都是之前linux的解决方案，在/etc/default下面创建一个rcS文件，在文件中配置UTC=no

但这个在最新的版本是不起作用，linux也是在不断改变。。。

后面在etc下面也找不到配置UTC=no的地方，这个就尴尬了。。。

因为网上很多都是使用timedatectl 命令设置的，因为是livecd，所以添加这个命令的地方不太方便，且比较容易出问题，因为系统可能在其他的地方设置是否采用UTC时间。

大概花了1-2小时，找到一个hwclock.sh 的文件，发现里面有一行shell 是判断 /etc/adjtime 文件里面是否开启UTC的，这个时候突然想到，可以使用timedatectl 命令设置 UTC = no， 发现在/etc/adjtime里面有配置，不过部分信息是二进制，无法查看，后面将adjtime 文件直接放入livecd中，发现这个问题解决了。



**注意: 上面所说的解决方式不能解决所有的问题，因为有的系统是默认BIOS时间也就是硬件时间是UTC，这个时候，我们如果在linux上改为LocalTime 时间就会出错，目前的方案主要针对windows系统，不能保证时间一定是准确的，如果需要时间准确最靠谱的办法还是获取网络时间。**