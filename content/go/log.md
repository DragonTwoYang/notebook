### go 日志
go 语言的基本库提供了丰富的日志功能，标准库中的log包为我们提供了可以满足日常使用的基本功能。


日志功能是一个健壮的程序所必备的模块。如果某个程序没有日志模块， 一旦程序出错， debug的代价可能会很大。
在产品的运营过程中，我们可能会看到这种情况， 软件在我们测试环节是ok的， 但到用户手中就出现不稳定。
特别是某些软件依赖的环境相对复杂，出现的概率更高。
所以，一个完善的日志模块是健壮程序所必须的。

```go
//   Ltime                         // 时间：01:23:23
//   Lmicroseconds                 // 微秒分辨率：01:23:23.123123（用于增强Ltime位）
//   Llongfile                     // 文件全路径名+行号： /a/b/c/d.go:23
//   Lshortfile                    // 文件无路径名+行号：d.go:23（会覆盖掉Llongfile）
//   LstdFlags     = Ldate | Ltime // 标准logger的初始值
package main

import (
    "log"
    "os"
)

func main() {
    file, err := os.
}
```


