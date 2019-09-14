#### I/O 多路复用之select

linux 下 I/O模型主要有以下几种



- 阻塞I/O模型
- 非阻塞I/O模型
- I/O 复用
- 信号驱动异步模型
- 异步I/O模型



这里先简单介绍下上面的几种模型：

阻塞I/O模型：当我们调用系统I/O接口去读取数据时，如果数据还没有准备好(或说数据为空，例如等待网络响应)，这个时候，如果数据没有准备好，调用会一直阻塞，直到数据到来才会重新唤醒线程。例如：recoverform



非阻塞I/O模型：当我们调用系统I/O接口去读取数据时，如果数据还没有准备好(或说数据为空，例如等待网络响应)，这个时候内核会直接返回一个错误，并不会阻塞线程



I/O 多路复用：上面2种模型对系统的资源利用率都不高所以有了I/O多路复用的技术，linux 支持的多路复用技术主要有select， poll， 和epoll三种(主流)。



信号驱动模型：首先开启套接口信号驱动I/O功能, 并通过系统调用sigaction安装一个信号处理函数（此系统调用立即返回，进程继续工作，它是非阻塞的）。当数据报准备好被读时，就为该进程生成一个SIGIO信号。随即可以在信号处理程序中调用recvfrom来读数据报，井通知主循环数据已准备好被处理中。也可以通知主循环，让它来读数据报。



异步I/O模型：告知内核启动某个操作，并让内核在整个操作完成后(包括将数据从内核拷贝到用户自己的缓冲区)通知我们。这种模型与信号驱动模型的主要区别是：信号驱动I/O：由内核通知我们何时可以启动一个I/O操作；异步I/O模型：由内核通知我们I/O操作何时完成。



##### select函数调用

select函数介绍：
 int select(int maxfd, fd_set *readset, fd_set *writeset, fd_set *exceptset, const struct timeval *timeout);
 功能：轮询扫描多个描述符中的任一描述符是否发生响应，或经过一段时间后唤醒

| 参数      | 名称                       | 说明                   |
| --------- | -------------------------- | ---------------------- |
| maxfd     | 指定要检查文件描述符的范围 | 所检测描述符的最大值+1 |
| readset   | 可读文件描述符             | 监测可读文件描述符     |
| writeset  | 可写文件描述符             | 监测可写文件描述符     |
| exceptset | 异常文件描述符             | 监测异常文件描述符     |
| timeout   | 超时                       | 超过规定时间后唤醒     |

```c
将select函数的timeout参数设置为NULL则永远等待  
```

```c
/初始化描述符集  
void FD_ZERO(fd_set *fdset);  
  
//将一个描述符添加到描述符集  
void FD_SET(int fd, fd_set *fdset);  
  
//将一个描述符从描述符集中删除  
void FD_CLR(int fd, fd_set *fdset);  
  
//检测指定的描述符是否有事件发生  
int FD_ISSET(int fd, fd_set *fdset);  
```



```c
#include <string.h>  
#include <stdio.h>  
#include <stdlib.h>  
#include <unistd.h>  
#include <sys/select.h>  
#include <sys/time.h>  
#include <sys/socket.h>  
#include <netinet/in.h>  
#include <arpa/inet.h>  
  
int main(int argc,char *argv[])  
{  
    int udpfd = 0;  
    struct sockaddr_in saddr;  
    struct sockaddr_in caddr;  
  
    bzero(&saddr,sizeof(saddr));  
    saddr.sin_family = AF_INET;  
    saddr.sin_port   = htons(8000);  
    saddr.sin_addr.s_addr = htonl(INADDR_ANY);  
      
    bzero(&caddr,sizeof(caddr));  
    caddr.sin_family  = AF_INET;  
    caddr.sin_port    = htons(8000);  
      
    //创建套接字  
    if( (udpfd = socket(AF_INET,SOCK_DGRAM, 0)) < 0)  
    {  
        perror("socket error");  
        exit(-1);  
    }  
      
    //套接字端口绑字  
    if(bind(udpfd, (struct sockaddr*)&saddr, sizeof(saddr)) != 0)  
    {  
        perror("bind error");  
        close(udpfd);         
        exit(-1);  
    }  
  
    printf("input: \"sayto 192.168.220.X\" to sendmsg to somebody\033[32m\n");    
    while(1)  
    {     
        char buf[100]="";     
        fd_set rset;    //创建文件描述符的聚合变量    
        FD_ZERO(&rset); //文件描述符聚合变量清0  
        FD_SET(0, &rset);//将标准输入添加到文件描述符聚合变量中  
        FD_SET(udpfd, &rset);//将udpfd添加到文件描述符聚合变量中        
        write(1,"UdpQQ:",6);  
          
        if(select(udpfd + 1, &rset, NULL, NULL, NULL) > 0)  
        {  
            if(FD_ISSET(0, &rset))//测试0是否可读写  
            {                 
                fgets(buf, sizeof(buf), stdin);  
                buf[strlen(buf) - 1] = '\0';  
                if(strncmp(buf, "sayto", 5) == 0)  
                {  
                    char ipbuf[16] = "";  
                    inet_pton(AF_INET, buf+6, &caddr.sin_addr);//给addr套接字地址再赋值.  
                    printf("\rsay to %s\n",inet_ntop(AF_INET,&caddr.sin_addr,ipbuf,sizeof(ipbuf)));  
                    continue;  
                }  
                else if(strcmp(buf, "exit")==0)  
                {  
                    close(udpfd);  
                    exit(0);  
                }  
                sendto(udpfd, buf, strlen(buf),0,(struct sockaddr*)&caddr, sizeof(caddr));  
            }  
            if(FD_ISSET(udpfd, &rset))//测试udpfd是否可读写  
            {  
                struct sockaddr_in addr;  
                char ipbuf[INET_ADDRSTRLEN] = "";  
                socklen_t addrlen = sizeof(addr);  
                  
                bzero(&addr,sizeof(addr));  
                  
                recvfrom(udpfd, buf, 100, 0, (struct sockaddr*)&addr, &addrlen);  
                printf("\r\033[31m[%s]:\033[32m%s\n",inet_ntop(AF_INET,&addr.sin_addr,ipbuf,sizeof(ipbuf)),buf);  
            }  
        }  
    }  
      
    return 0;  
}  
```






```c
#include <sys/select.h>
#include <sys/time.h>
#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>

int main()
{
    fd_set rd;
    struct timeval tv;
    int err;
	FD_ZERO(&rd);
	FD_SET(0,&rd);

	tv.tv_sec = 5;
	tv.tv_usec = 0;
	err = select(1,&rd,NULL,NULL,&tv);

	if(err == 0) //超时
	{
    	printf("select time out!\n");
	}
	else if(err == -1)  //失败
	{
    	printf("fail to select!\n");
	}
	else  //成功
	{
    	printf("data is available!\n");
	}
	return 0;
}
```

