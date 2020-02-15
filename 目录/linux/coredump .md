#### core dump文件归档

##### 1.core dump文件是什么，如何生成的？

当程序运行的过程中异常终止或崩溃，操作系统会将程序当时的内存状态记录下来，保存在一个文件中，这种行为就叫做Core Dump（中文有的翻译成“核心转储”)。**core dump 对于编程人员诊断和调试程序是非常有帮助的**，因为对于有些程序错误是很难重现的，例如指针异常，而 core dump 文件可以再现程序出错时的情景。



目前，我们在uBeanOS桌面系统和服务器系统里均打开了core dump文件的生成开关。

获取内核转储（core dump）的最大好处就是，它能保存问题发生时的状态。只要有问题发生时程序的可执行文件和内核转储，就可以知道进程当时的状态。比如在不清楚bug复现的方法情况下，或者bug极其罕见，只要有内核转储文件，我们就可以调试。



##### 2.如何设置系统的coredump文件生成？

大多数linux发行版默认是关闭内核转储功能的。

使用ulimit命令可以查看内核转储是否有效：

```shell
ulimit -c
0
```

 如果输出是0， 表明内核转储功能无效， 可以使用下面命令开启：

```shell
ulimit -c unlimited
```

上面命令的意思是不限制核心转储文件的大小。



###### 2.1关闭uBeanos 的核心转储功能：

 修改/etc/security/limits.d/90-ubeanos-default-limits.conf 文件

 将 *   hard core   unlimited  这一行的ulimited 改成 0 ， 改完之后是： *  hard   core   0

 将 *   soft  core  unlimited  这一行的ulimited 改成  0， 修改之后是： *  sofr   core   0



###### 2.2 如何在专有的目录中生成内核转储文件

核心转储文件会默认生成在当前目录中，但在生成环境中，会希望可以将内核转储文件生成在固定的位置。

还有就是系统在生产环境下，生产的内核转储文件过大或过多，生成的文件会消耗磁盘空间，影响系统的正常运行，此时可能会采用分区单独存放转储文件的方案。

转储文件保存位置可以通过sysctl 变量 kernel.core_pattern 设置

kernel.core_pattern=/var/core/%t-%e.core

kernel.core_pattern 中可以设置的格式符：

| 格式符 | 说明                               |
| :----- | ---------------------------------- |
| %%     | %符号本身                          |
| %p     | 被转储进程的ID(PID)                |
| %u     | 被转储进程的真实用户ID（real UID） |
| %g     | 被转储进程的真实组ID（real GID）   |
| %s     | 引发转储的信号编号                 |
| %t     | 转储时间（UNIX时间秒数）           |
| %h     | 主机名                             |
| %e     | 可执行文件名称                     |
| %c     | 转储文件的大小限制                 |

*需要注意：指定生成文件的目录需要所有用户都有可执行权限，要不会生成失败*



###### 2.3压缩核心转储文件

刚刚的kernel.core_pattern 可以加入管道的功能

kernel.core_pattern=|/usr/local/bin/core_helper %t %e %p %c 

```shell
#!/bin/bash
#core_helper

exec gzip - > /var/core/$1-$2-$3-$4.core.gz
```



###### 2.4利用内核转储掩码排除共享内存

有的应用程序会使用多个进程，多个进程间可能会使用共享内存，核心转储文件会保存共享内存，当多个进程进行共享内存时候，只需要一个进程进行保存即可。

设置方法非常简单，通过/proc/<PID>/coredump_filter进行。

coredump_filter 使用比特掩码表示内存类型

| 比特掩码 | 内存类型             |
| :------- | -------------------- |
| bit 0    | 匿名专用内存         |
| bit 1    | 匿名共享内存         |
| bit 2    | file-backed 专用内存 |
| bit 3    | file-backed 共享内存 |
| bit4     | ELF 文件映射         |



##### 3.归档策略如何选择？

-  每天凌晨进行归档

- core dump文件进行压缩

- 保存7天的归档文件

- 保存最近的文件

- 可以设置备份目录的大小

  

简单实现代码如下：

```shell
#!/bin/bash 

#This script aim is to archive uBeanOS core dump file.

DIR=/var/core
COREDUMPS=/var/core/dumps
MAXSIZE=100000 #kb

#tar cvzf  coredumps.1.gz  COREDUMPS/core.*  

check(){
	files=$(ls $DIR/core.*)
	if [ X"$files" == X ];then
		exit 0
	fi
}

dir_size(){
	dir=$1
	size=$(du -s $dir | awk -F ' ' '{print $1}')
	echo $size
}

get_max_log(){
	maxnum=0
	files=`ls $COREDUMPS | sort`
	for i in $files; do
		num=$(echo $i | awk -F '.' '{print $2}')
		if [ $num -gt $maxnum ];then
			maxnum=$num
		fi
	done
	echo $maxnum
}

get_minx_log(){
	maxnum=100000
	files=`ls $COREDUMPS | sort`
	for i in $files; do
		num=$(echo $i | awk -F '.' '{print $2}')
		if [ $num -lt $maxnum ];then
			maxnum=$num
		fi
	done
	echo $maxnum
}

update_file_num(){
	max=$(get_max_log)
	min=$(get_minx_log)
	if [ $max == $min ];then
		n=$(expr $min + 1)
		if [ -f $COREDUMPS/coredumps.${min}.gz ];then
			mv $COREDUMPS/coredumps.${max}.gz  $COREDUMPS/coredumps.${n}.gz
		fi
		return 

	fi	
	for i in `seq $max -1 $min`
	do
		n=$(expr $i + 1)
		if [ -f $COREDUMPS/coredumps.${i}.gz ];then
			mv $COREDUMPS/coredumps.${i}.gz $COREDUMPS/coredumps.${n}.gz
		fi
	done
}

check_dir_size(){
	size=$(dir_size $COREDUMPS)
	if [ ${size} -gt ${MAXSIZE} ]; then
		num=$(get_max_log)
		rm $COREDUMPS/coredumps.${num}.gz
		check_dir_size
	fi
}

archive_files(){
	files=$(ls ${DIR}/core.*)
	if [ X"$files" != X ];then 
		tar Pcvzf  $COREDUMPS/coredumps.1.gz  ${DIR}/core.*  
	fi
}

check
get_max_log
get_minx_log
check_dir_size
update_file_num
archive_files

```



补充：crontab 设置定时任务：

/var/spool/cron/crontabs/root (0600)

root 表示是root用户创建的任务

\* * * * * /bin/coreDumpTar.sh   //每分钟执行一次coreDumpTar.sh 脚本