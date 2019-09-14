# linux kernel编译

## linux kernel 简介

*操作系统是一个用来和硬件打交道并为用户程序提供一个有限服务集的低级支撑软件。</br> 一个计算机系统是一个硬件和软件的共生体，它们互相依赖，不可分割。</br>计算机的硬件，含有外围设备、处理器、内存、硬盘和其他的电子设备组成计算机的发动机。</br>但是没有软件来操作和控制它，自身是不能工作的。完成这个控制工作的软件就称为操作系统，在Linux的术语中被称为“内核”，也可以称为“核心”。</br>Linux内核的主要模块（或组件）分以下几个部分：存储管理、CPU和进程管理、文件系统、设备管理和驱动、网络通信，以及系统的初始化（引导）、系统调用等。*

上面是百度百科上关于操作系统的介绍，内核(kernel)是一个操作系统的核心。</br>
linux 内核和其他的内核很大一个不同之处是linux 内核源代码是可以查看且编译的。</br>
linxu 内核支持很多的平台和硬件， 不同的平台编译上可能会有细微差异，下面主要讲述x86平台的内核编译。

## kernel 源码获取

  你可以从[kernel achive](https://www.kernel.org/) 下载你所指定的内核版本。</br>如果是使用， 可以使用stable分支的， 如果仅仅是学习和试验，版本差异不明显。</br>除了从官网之外， 你还可以从github上下载kernel的源码。

## kernel 编译
 
### kernel config

    make menuconfig

### kernel build

    make build

### kernel module build

    make modules

### kernel 安装

    make install 

### kernel module 安装

    make modules install
