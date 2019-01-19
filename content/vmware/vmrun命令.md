## vmrun 命令使用笔记
1.vmrun是vmware发布用来控制vmware虚拟机的外部命令，与之类似功能的还有一个vix的库。
vix库提供了和vmrun命令能力相当的接口。凡是使用vmrun可以实现的功能，使用vix库也可以实现。vix的主要优势在于可以更获取更多虚拟机的状态。例如，如果使用vix启动一个虚拟机，我们可以获取虚拟机启动过程中的状态，而vmrun只能获取命令执行成功还是失败。
因为vmrun使用简单可靠，本次笔记主要记录vmrun的常用用法。

### vmrun 有哪些用法
通过vmrun --help ，我们可以查看vmrun的主要用法和简单使用实例。

### vmrun 常用命令
启动虚拟机
vmrun.exe -T ws start "xxx.vmx" 

后台启动虚拟机
vmrun.exe -T ws start "xxx.vmx" nogui

关闭虚拟机(正常关闭)
vmrun.exe -T ws stop "xxx.vmx" soft

关闭虚拟机(强制关闭)
vmrun.exe -T ws stop "xxx.vmx" hard

创建快照
vmrun.exe -T ws snapshot "xxx.vmx" mySanpshot

删除快照
vmrun.exe -T ws deleteSnapshot "xxx.vmx" mySanpshot

创建共享目录
#TODO

执行程序
vmrun.exe -T ws -gu root -gp 1 runProgramInGuest  "E:\虚拟机\Fedora 27.vmx"  "/bin/bash" "/home/yan/1.sh" "1" "2"

确定Guest os中文件是否存在
vmrun.exe -T ws -gu root -gp 1  fileExistsInGuest   "E:\虚拟机\Fedora 27.vmx"  "/1.log"

从guset拷贝文件到host
vmrun.exe -T ws -gu root -gp 1 CopyFileFromGuestToHost "E:\虚拟机\Fedora 27.vmx"  "/1.log"  E:\testCopy.txt