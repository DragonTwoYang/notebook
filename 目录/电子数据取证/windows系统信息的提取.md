### windows 系统信息的提取

   取证方向有一项非常重要的信息提取，就是提取操作系统的基本信息，以帮助办案人员查看嫌疑人的电脑相关信息。不同的操作系统os基本信息提取方式差别较大，今天就单独说明windows下如何提取系统信息。

windows上非常多的信息存储在注册表里面，windows注册表的解析也是解析windows基本信息必不可少的基本工作，下面我们不会讲解如何解析windows的注册表，仅仅是说明信息存储在注册表的位置。



windows os基本信息主要包含以下几部分：

- 操作系统基本信息(名称，系统类型，最后一次关机时间，操作系统版本， 当前build版本，产品ID， 安装时间，系统根路径等)
- 共享信息
- 系统安装软件(windows自带，用户自行安装)
- 应用分类
- 快捷方式
- 打印机
- 服务信息
- 用户信息
- 系统开关机时间
- 网络配置
- 自启动程序
- 默认应用



##### 操作系统基本信息

操作系统的基本信息基本都可以从注册表中获取

 注册表位置：计算机\HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion

这个key下面有非常多操作系统的基本信息

ProductName  = Windows 10 Home China   //用户产品名称，即操作系统的名称

CurrentVersion = 6.3  //当前OS的基本版本

InstallDate = 1527130024  //系统安装日期

InstallTime = 131716036241838018  //安装时间

PathName = C:\Windows  //系统根目录

SoftwareType = System  //系统软件(WINDOWS NT)

ProductId = 00342-34765-50054-AAOEM  //产品ID



注册表位置： 计算机\HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\Control\Windows

这个key下面有系统最后一次关机时间

ShutdownTime  //这个里面保存了最后一次关机时间

注册表位置：计算机\HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\Control\Session Manager\Environment

PROCESSOR_ARCHITECTURE = AMD64  //确定操作系统类型，32位还是64位



##### 共享信息

这里共享信息就是单指windows系统把目录共享给局域网内其他主机使用的信息。

文件夹共享信息也被记录进注册表，我们可以通过注册表来获取共享文件的信息。

注册表位置：计算机\HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\LanmanServer\Shares

共享文件的名称就会在上面注册表下。

key： ubuntu16.04 



value:

CATimeout=0
CSCFlags=2096
MaxUses=9
Path=E:\ubuntu16.04
Permissions=0
Remark=this is a test
ShareName=ubuntu16.04
Type=0

简单说明下，MaxUses: 最大可以连接数目   Remark：描述   ShareName：共享名称   Path：共享路径



##### 快捷方式解析

windows 上面有非常多的快捷方式，比如桌面上存放的某个程序的打开快捷方式。

windows上面的快捷方式是直接采用二进制存储的，有自己定义的格式(相对省空间，可以暂用更少的磁盘)。

linux上面的快捷方式也有自身定义的格式，但是是文本存储方式，用户可以简单的修改和查看(简单明了)。

快捷方式的解析，请查看另外的章节。



##### 系统开关机时间

系统的开关机时间是一个非常有价值的信息，可以判断用户使用系统的时间。

系统的开关机时间并没用直接存储，注册表中有最后一次的关机时间，但仅仅是最后一次。

注册表中虽然没用开关机的时间，但我们可以通过时间日志来获取系统的开关机时间。

windows下面的事件日志也就是eventlog， windows的事件日志比linux强上不少，因为windows的事件日志不仅仅给不同的事件分配ID，大家可以通过事件ID来区分事件类型，还有windows为app开发也预留了便于使用的接口，可以把app的事件日志和系统事件日志统一起来，便于用户查看。

系统的开关机时间，主要使用了2个事件：

开机：6009

关机：6006 

通过以上2个事件，我们可以很方便的查找出系统的开关机时间。



##### 用户信息



##### 