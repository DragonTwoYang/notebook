# 定制linux系统

## 获取系统镜像

  原始iso可以从开源镜像站获取(镜像齐全),或者从不同发行版的官网下载。</br>
  [163开源镜像地址](https://mirrors.163.com/) </br>
  [阿里云开源镜像](https://opsx.alibaba.com/mirror)

## 镜像格式说明

iso,国际标准光盘文件系统格式。符合ISO 9660标准的光盘镜像文件格式,文件扩展名通常为iso。</br>
这种文件可以简单的理解为复制光盘上全部信息而形成的镜像文件。

## linux手动修改ubuntu镜像

*以下以ubuntu.iso 代表待修改ISO镜像名字，以下以/home/yan为工作目录*</br>

### 1.挂载&拷贝ISO

先挂在ISO镜像，再在工作目录创建livecd并拷贝存放ISO的内容

    #sudo  mkdir  /home/yan/livecd
    #sudo  mount  -o  loop  ubuntu.iso   /mnt
    #sudo  rsync  -a  /mnt  /home/yan/livecd

### 2.解压缩squashfs Image

解压缩ISO里的squashfs Image 文件，为后续修改做准备。

    #sudo  apt-get   install   squashfs-tools
    #sudo  mv  /home/yan/livecd/casper/filesystem.squashfs/home/cdos
    #sudo  unsquashfs  /home/yan/filesystem.squashfs  #会在/home/yan/创建squashfs-tools文件夹

### 3.挂载相关文件系统

    #sudo  mount  -t  proc  proc  /home/yan/squashfs-tools/proc
    #sudo  mount  -t  sysfs  sysfs  /home/yan/squashfs-tools/sys
    #sudo  mount  -t  devtmpfs  devtmpfs/home/yan/squashfs-tools/dev

### 4.需改ISO内容

*在对ISO做修改前，需要先拷贝相关内容到/home/cdos/squashfs-tools/目录下，然后chroot后进行修改*</br>

    #sudo  cp  xxx_1.0.0_amd64.deb  /home/yan/squashfs-tools/tmp/
    #sudo  chroot  /home/yan/squashfs-tools/
    #dpkg   -i  /tmp/xxx_1.0.0_amd64.deb
    #dpkg -l | grep ii | awk '{print $2,$3}' > filesystem.manifest

### 5.重构squashfs Image

    #exit      #退出chroot环境
    #sudo  umount  -t  /home/yan/squashfs-tools/proc
    #sudo  umount  -t   /home/yan/squashfs-tools/sys
    #sudo  umount  -t   /home/yan/squashfs-tools/dev
    #sudo  mv  /home/yan/squashfs-tools/filesystem.manifest /home/yan/livecd/casper/filesystem.manifest
    #sudo  mksquashfs  /home/yan/squashfs-tools   /home/yan/livecd/casper/filesystem.squashfs

### 6.更新md5sum.txt

    #cd /home/yan/livecd/
    #find -type f -print0 | sudo xargs -0 md5sum | grep -v ./md5sum.txt |sudo  tee  md5sum.txt

### 7.重新制作iso镜像

    #sudo  apt-get  install  mkisofs
    #sudo  cd  /home/yan/livecd/
    #sudo mkisofs -D -r -V "elef" -cache-inodes -J -l -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table  -o ../yan_v001.iso .

## 制作流程简述

**通过上面的修改流程， 我们可以轻易的发现， 修改iso的流程其实很简单。**</br>

 1. 先解压iso镜像  
 2. 再解压开里面系统的文件  
 3. 将解压的系统挂载到宿主机器上 
 4. 修改
 5. 逆向压缩回来
