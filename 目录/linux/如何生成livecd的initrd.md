##### 如何生成livecd的initrd文件

###### 一.initrd作用概述

linux kernel 启动之后，便会加载initrd。initrd其实就是一个简单的rootfs， 我们可以自己通过busybox定制一个。

initrd脚本执行过程中可以做各种事情：

1.把系统安装器放入initrd中， 例如debian-installer程序

2.将部分驱动的加载放入initrd中， 这样可以使得编译的kernel更小，更精简。举个简单例子，磁盘的驱动就可以放入initrd中， 这样kernel中就不用编译所以的磁盘驱动程序。

3.还有一个就是大家常见的livecd， 也是在initrd中只读挂载rootfs。

当然，initrd中可以定制各种你的需求， 不仅限于上面所说。



###### 二. livecd 实现原理简述

Ubuntu上实现livecd 主要依靠以下几个组件：

1.linux kernel 

2.initrd 

3.完整系统的rootfs(squashfs格式)

4.每次进入都重置(mount ro + overly)

简单说下基本原理：

系统kernel启动之后， initrd会被加载进入内存中， kernel会将程序的执行流程由kernel 转向 initrd中的init程序。

init程序会casper程序的脚本， 主要就是环境的设置，驱动的加载， 查找rootfs文件并以只读的形式挂载进入系统。

系统挂载之后，initrd程序会执行rootfs的init程序， 如何是安装systemd的系统，一般就是执行systemd程序，也就是大家说熟悉的1号进程。systemd会启动服务，进入桌面等。



###### 三，如何生成livecd启动所需要initrd

ubuntu下，livecd中initrd中的脚本都在casper程序下面，你只需要定制好自己系统的rootfs， 再安装casper包就可以通过update-initramfs命令生成你需要的initrd文件。

debian下面， 你可以通过live-boot和live-build包，直接生成一个完整可用的livecd iso， 里面包含了你所需要的initrd等文件。



###### 四，Livecd如何从硬盘启动

```shell
sudo vim /etc/grub.d/40_custom
```

```shell
menuentry "ubuntu 18.04 ISO" {
set isofile="/home/yan/ubuntu-18.04-desktop-amd64.iso"
loopback loop (hd0,1)$isofile
echo "Starting $isofile..."
linux (loop)/casper/vmlinuz boot=casper iso-scan/filename=${isofile} quiet splash
initrd (loop)/casper/initrd.lz
}
```



###### 五，如何从硬盘启动自定义的rootfs

新的需求，如何从硬盘中启动我们自己定义的rootfs文件(里面是一个完整的系统，包含initd和kernel)。

目前Ubuntu上面处理livecd启动流程的是casper程序，这个包的上游是debian的live-boot包。

还是debian的包名更为合理和易懂点。

第六节将会介绍live-boot中，我们可以传递怎么样的参数给live-boot程序。

1.将filesystem.squashfs 放入/casper目录下面

2.在/etc/grub.d/40_custom 文件中加入下面的代码

```shell
menuentry 'Live cd' {
    set imgfile=/casper/filesystem.squashfs
    loopback loop0 $imgfile
    linux (loop0)/vmlinuz root=$loop0 boot=casper looptype=squashfs debug=1 udev cdroot
    initrd (loop0)/initrd.img
}
```

3.生成新的grub配置文件

```shell
sudo update-grub2
```

4.重启电脑，grub菜单中选择live cd选项



###### 六, debian live-boot程序参数简介

```
live-boot - System Boot Scripts
```

```
live-boot is only activated if 'boot=live' was used as a kernel parameter.

```

```tiki wiki
  live-boot currently features the following parameters.

access=ACCESS
		Set  the accessibility level for physically or visually impaired users. ACCESS must be one of v1, v2,  v3,  m1,  or  m2.  v1=lesser  visual  impairment,  v2=moderate  visual impairment, v3=blindness, m1=minor motor difficulties, m2=moderate motor difficulties.

console=TTY,SPEED
		Set the default  console  to  be  used  with  the  "live-getty"  option. Example:"console=ttyS0,115200"

debug
		Makes initramfs boot process more verbose.
        Use: debug=1
        Without setting debug to a value the messages may not be shown.

fetch=URL

httpfs=URL
   Another form of netboot by downloading a squashfs image from a given url.   The  fetch  method  copies  the  image to ram and the httpfs method uses fuse and httpfs2 to mount the image in place. Copying to ram requires more memory and might take a long time for large images. However, it is more likely to work correctly because it does not require
networking afterwards and the system operates faster once booted because it  does  not require to contact the server anymore.Due  to  current  limitations  in  busybox's  wget  and DNS resolution, an URL can not contain a hostname but an IP only.
           Not working: http://example.com/path/to/your_filesystem.squashfs
           Working: http://1.2.3.4/path/to/your_filesystem.squashfs
           Also note that therefore it's  currently  not  possible  to  fetch  an  image  from  a namebased  virtualhost  of  an  httpd  if  it  is  sharing  the ip with the main httpd instance. You may also use the live iso image in place of the squashfs image.
           
iscsi=server-ip[,server-port];target-name
	Boot from an iSCSI target that has an iso or disk live image as one of its  LUNs.  The specified  target  is searched for a LUN which looks like a live media. If you use the iscsitarget software iSCSI target solution which is packaged in Debian your  ietd.conf
           might look like this:
           # The target-name you specify in the iscsi= parameter
           Target <target-name>
             Lun 0 Path=<path-to-your-live-image.iso>,Type=fileio,IOMode=ro
             #  If  you  want  to  boot  multiple  machines you might want to look at 		         tuning some
           parameters like
             # Wthreads or MaxConnections

findiso=/PATH/TO/IMAGE
		Look for the specified ISO file on all disks where it usually looks for the  .squashfs file (so you don't have to know the device name as in fromiso=....).

fromiso=/PATH/TO/IMAGE
		Allows to use a filesystem from within an iso image that's available on live-media.

ignore_uuid
		Do  not  check  that any UUID embedded in the initramfs matches the discovered medium.live-boot may be told to generate a UUID by setting LIVE_GENERATE_UUID=1 when building the initramfs.

verify-checksums
		If  specified,  an MD5 sum is calculated on the live media during boot and compared to the value found in md5sum.txt found in the root directory of the live media.

ip=[DEVICE]:[CLIENT_IP]:[NETMASK]:[GATEWAY_IP]:[NAMESERVER]
       [,[DEVICE]:[CLIENT_IP]:[NETMASK]:[GATEWAY_IP]:[NAMESERVER]]
		Let  you  specify  the  name(s)  and  the  options  of the interface(s) that should be configured at boot time. Do not specify this if you want to  use  dhcp(default).  It will be changed in a future release to mimick official kernel boot param specification(e.g. ip=10.0.0.1::10.0.0.254:255.255.255.0::eth0,:::::eth1:dhcp).

ip=[frommedia]
		If this variable is set, dhcp and static configuration are just skipped and the system will use the (must be) media-preconfigured /etc/network/interfaces instead.

{live-media|bootfrom}=DEVICE
		If you specify one of this two equivalent forms, live-boot will first try to find this device for the "/live" directory where the read-only root filesystem should reside. If it did not find something usable, the normal scan for block devices is performed.
Instead  of  specifing  an  actual device name, the keyword 'removable' can be used to
limit the search of acceptable live media to removable type only.  Note  that  if  you
want  to  further  restrict  the  media  to  usb  mass  storage  only, you can use the
'removable-usb' keyword.

{live-media-encryption|encryption}=TYPE
		live-boot will mount the encrypted rootfs TYPE, asking the passphrase, useful to build paranoid  live  systems  :-).  TYPE supported so far are "aes" for loop-aes encryption type.

live-media-offset=BYTES
		This way you could tell live-boot that your image starts at offset BYTES in the  above specified  or  autodiscovered device, this could be useful to hide the Debian Live iso or image inside another iso or image, to create "clean" images.

live-media-path=PATH
		Sets the path to the live filesystem on the medium. By default, it is set  to  '/live' and you should not change that unless you have customized your media accordingly.

live-media-timeout=SECONDS
		Set  the  timeout in seconds for the device specified by "live-media=" to become ready before giving up.

module=NAME
		Instead of using the default optional file  "filesystem.module"  (see  below)  another file  could  be  specified  without  the  extension  ".module"; it should be placed on "/live" directory of the live medium.

netboot[=nfs|cifs]
		This tells live-boot to perform  a  network  mount.  The  parameter  "nfsroot="  (with optional  "nfsopts="),  should  specify  where is the location of the root filesystem.With no args, will try cifs first, and if it fails nfs.

nfsopts=
		This lets you specify custom nfs options.

nofastboot
		This parameter disables the default disabling of filesystem checks in  /etc/fstab.  If you  have  static filesystems on your harddisk and you want them to be checked at boot time, use this parameter, otherwise they are skipped.

nopersistence
		disables the "persistence" feature, useful if the bootloader (like syslinux) has  been installed with persistence enabled.noeject Do not prompt to eject the live medium.

ramdisk-size
		This  parameters  allows  to  set  a custom ramdisk size (it's the '-o size' option of tmpfs mount). By default, there is no ramdisk  size  set,  so  the  default  of  mount applies  (currently  50%  of available RAM). Note that this option has no currently no effect when booting with toram.

swapon
       This parameter enables usage of local swap partitions.

persistence
		live-boot will probe devices for persistence media. These can be partitions (with  the correct  GPT  name),  filesystems  (with  the  correct label) or image files (with the correct    file    name).    Overlays    are    labeled/named    "persistence"    (seepersistence.conf(5)).  Overlay  image  files  have  extensions  which determines their filesystem, e.g. "persistence.ext4".

persistence-encryption=TYPE1,TYPE2 ... TYPEn
		This option determines which types of encryption that we allow to be used when probingdevices  for  persistence media. If "none" is in the list, we allow unencrypted media;if "luks" is in the list, we allow LUKS-encrypted media. Whenever a device  containing encrypted  media  is  probed the user will be prompted for the passphrase. The default  value is "none".

persistence-media={removable|removable-usb}
		If you specify the  keyword  'removable',  live-boot  will  try  to  find  persistence partitions  on  removable  media  only.  Note that if you want to further restrict the media to usb mass storage only, you can use the 'removable-usb' keyword.

persistence-method=TYPE1,TYPE2 ... TYPEn
    	This option determines which types of persistence media we allow. If "overlay"  is  in the  list,  we  consider  overlays  (i.e.  "live-rw"  and  "home-rw").  The default is
"overlay".

persistence-path=PATH
    	live-boot will look for persistency files in the root directory of a  partition,  with this  parameter,  the path can be configured so that you can have multiple directories on the same partition to store persistency files.

persistence-read-only
		Filesystem changes are not saved back to persistence media.  In  particular,  overlays and netboot NFS mounts are mounted read-only.

persistence-storage=TYPE1,TYPE2 ... TYPEn
		This option determines which types of persistence storage to consider when probing for persistence media. If "filesystem" is in the list, filesystems  with  matching  labels will  be  used;  if "file" is in the list, all filesystems will be probed for archives and image files with matching filenames. The default is "file,filesystem".

persistence-label=LABEL
		live-boot will use the name  "LABEL"  instead  of  "persistence"  when  searching  for persistent  storage.  LABEL  can  be any valid filename, partition label, or GPT name.  This option replaces the less flexible persistent-subtext option  from  version  2  of live-boot.  If you wish to continue using legacy names for persistent storage, use the full name with this option, e.g. persistence-label=live-rw-foo

       quickreboot
           This option causes live-boot to reboot without  attempting  to  eject  the  media  and
           without asking the user to remove the boot media.

showmounts
		This  parameter  will  make  live-boot  to  show  on  "/"  the  ro filesystems (mostly  compressed) on "/lib/live". This is not enabled  by  default  because  could  lead  to  problems by applications like "mono" which store binary paths on installation.
silent If you boot with the normal quiet parameter, live-boot hides most messages of its own. When adding silent, it hides all.

todisk=DEVICE
		Adding this parameter, live-boot will try to copy the entire read-only  media  to  the specified  device before mounting the root filesystem. It probably needs a lot of free space.  Subsequent  boots  should  then  skip  this  step   and   just   specify   the "live-media=DEVICE" boot parameter with the same DEVICE used this time.

toram
		Adding  this  parameter,  live-boot  will try to copy the whole read-only media to the computer's RAM before mounting the root filesystem. This could  need  a  lot  of  ram,according to the space used by the read-only media.

union=aufs|unionfs
		By default, live-boot uses aufs. With this parameter, you can switch to unionfs.
```

