### 修改windows密码工具之chtpw介绍说明

------

### 一.概述

chntpw是一个修改windows密码的工具，可以在linux上用，也可以在windows上使用。

chntpw是几个单词的缩写，ch：change  nt：windows NT  pw：password  

chntpw是通过将windows SAM文件中密码的hash长度制成0来达到清除密码的作用。

chntpw 目前在ubuntu上是1.1最新版本，支持windows xp， server， window 7，8 ，10的密码修改。

------

### 二.sampasswd命令使用方式

  修改windows密码必须离线启动，其实就是启动一个windows PE或 linux LIVECD(离线方式可以修改文件，否则文件时无法写入的)。

我们以ubuntu 为例，介绍chngpw的工具使用：

- 先安装chngpw 工具：

```shell
sudo apt-get install chntpw  
```



- mount windows 系统所在分区(假定为sda1)

```shell
sudo mount /dev/sda1 /mnt 
```



- 进入SAM文件所在目录

  ```shell
  cd /mnt/Windows/System32/config
  ```

  

-   查看SAM文件中有哪些用户名

  ```shell
  sampasswd  -l  SAM
  ```



- 通过命令重置账户密码

  ```shell
  sampasswd -r -u admin  SAM
  ```



-   卸载挂载的windows 系统分区

  ```shell
  cd ~；sudo  unmount  /mnt
  ```

- 重启机器

  ```shell
  reboot
  ```

  

chntpw 命令提供了chntpw和 sampasswd2个命令行工具

如果时用户直接操作，chntpw命令提供menu菜单，但不利于脚本实现密码自动修改

除了chntpw命令，此包还有提供sampasswd命令，作用和chntpw一致，但这个命令非常方便脚本自动化执行。

通过chngpw -l 命令和sampasswd 命令可以很直观查看此工具有提供哪些功能。

------

### 三. chntpw 提供了哪些功能

修改密码只是chntpw中提供的一项功能，下面简单介绍下其他功能

- 查看用户的用户组
- 添加新的账号
- 查看用户ID
- 修改用户所属组
- 查看和导出已经修改注册表的工具

------

四. help

#sampasswd --hep

```
sampasswd version 0.2 140201, (c) Petter N Hagen
sampasswd  [-r|-l] [-H] -u <user> <samhive>
Reset password or list users in SAM database
Mode:
   -r = reset users password
   -l = list users in sam
Parameters:
   <user> can be given as a username or a RID in hex with 0x in front
   Example:
   -r -u theboss -> resets password of user named 'theboss' if found
   -r -u 0x3ea -> resets password for user with RID 0x3ea (hex)
   -r -a -> Reset password of all users in administrators group (0x220)
   -r -f -> Reset password of admin user with lowest RID
            not counting built-in admin (0x1f4) unless it is the only admin
   Usernames with international characters usually fails to be found,
   please use RID number instead
   If success, there will be no output, and exit code is 0
Options:
   -H : For list: Human readable listing (default is parsable table)
   -H : For reset: Will output confirmation message if success
   -N : No allocate mode, only allow edit of existing values with same size
   -E : No expand mode, do not expand hive file (safe mode)
   -t : Debug trace of allocated blocks
   -v : Some more verbose messages/debug

```



#chntpw --help

```
chntpw: change password of a user in a Windows SAM file,
or invoke registry editor. Should handle both 32 and 64 bit windows and
all version from NT3.x to Win8.1
chntpw [OPTIONS] <samfile> [systemfile] [securityfile] [otherreghive] [...]
 -h          This message
 -u <user>   Username or RID (0x3e9 for example) to interactively edit
 -l          list all users in SAM file and exit
 -i          Interactive Menu system
 -e          Registry editor. Now with full write support!
 -d          Enter buffer debugger instead (hex editor), 
 -v          Be a little more verbose (for debuging)
 -L          For scripts, write names of changed files to /tmp/changed
 -N          No allocation mode. Only same length overwrites possible (very safe mode)
 -E          No expand mode, do not expand hive file (safe mode)

Usernames can be given as name or RID (in hex with 0x first)

See readme file on how to get to the registry files, and what they are.
Source/binary freely distributable under GPL v2 license. See README for details.
NOTE: This program is somewhat hackish! You are on your own!

```

