**golong面试题**

1.关于defer的执行顺序问题

```go
package main

import (
    "fmt"
)

func main() {
    defer_call()
}

func defer_call() {
    defer func() { fmt.Println("打印前") }()
    defer func() { fmt.Println("打印中") }()
    defer func() { fmt.Println("打印后") }()

    panic("触发异常")
}
```

解答：

go中的defer 是后入先出的，压栈模式

程序发出panic之后，会在当前携程里面往外执行defer语句， 最后处理panic。



2.go语言中foreach问题

``

```go
type student struct {
    Name string
    Age  int
}

func pase_student() {
    m := make(map[string]*student)
    stus := []student{
        {Name: "zhou", Age: 24},
        {Name: "li", Age: 23},
        {Name: "wang", Age: 22},
    }
    for _, stu := range stus {
        m[stu.Name] = &stu
    }

}
```

go中的for， range关键字循环和java中的foreach一样，使用的是副本拷贝。

stu的指针所指向的值不同，但指针是一致的。最后都为wang，age 22.



3.go中闭包问题以及协程执行顺序问题

``

```go
func main() {
    runtime.GOMAXPROCS(1)
    wg := sync.WaitGroup{}
    wg.Add(20)
    for i := 0; i < 10; i++ {
        go func() {
            fmt.Println("A: ", i)
            wg.Done()
        }()
    }
    for i := 0; i < 10; i++ {
        go func(i int) {
            fmt.Println("B: ", i)
            wg.Done()
        }(i)
    }
    wg.Wait()
}
```



4.go组合继承问题

```go
type People struct{}

func (p *People) ShowA() {
    fmt.Println("showA")
    p.ShowB()
}

func (p *People) ShowB() {
    fmt.Println("showB")
}

type Teacher struct {
    People
}

func (t *Teacher) ShowB() {
    fmt.Println("teacher showB")
}

func main() {
    t := Teacher{}
    t.ShowA()
}
```

解答：输出应该是：showA，showB

这是Golang的组合模式，可以实现OOP的继承。 被组合的类型People所包含的方法虽然升级成了外部类型Teacher这个组合类型的方法（一定要是匿名字段），但它们的方法(ShowA())调用时接受者并没有发生变化。 此时People类型并不知道自己会被什么类型组合，当然也就无法调用方法时去使用未知的组合者Teacher类型的功能。



5.go关键字select的随机性

```go
func main() {
    runtime.GOMAXPROCS(1)
    int_chan := make(chan int, 1)
    string_chan := make(chan string, 1)
    int_chan <- 1
    string_chan <- "hello"
    select {
    case value := <-int_chan:
        fmt.Println(value)
    case value := <-string_chan:
        panic(value)
    }
}
```

解答：go中的select先收到哪个值是不确定的。

注意：如果string_chan := make(chan string），将会发送死锁。



6.defer的调用时机

```go
func calc(index string, a, b int) int {
    ret := a + b
    fmt.Println(index, a, b, ret)
    return ret
}

func main() {
    a := 1
    b := 2
    defer calc("1", a, calc("10", a, b))
    a = 0
    defer calc("2", a, calc("20", a, b))
    b = 1
}
```

解答：defer 函数里面的参数会先执行

输出：

​10 1 2 3
20 0 2 2
2 0 2 2
1 1 3 4



7.make的默认值问题

```go
func main() {
    s := make([]int, 5)
    s = append(s, 1, 2, 3)
    fmt.Println(s)
}
```

解答：输出是0，0，0，0，0，1，2，3

因为make初始化是5个。



8.go中常量问题

```go
package main
const cl  = 100

var bl    = 123

func main()  {
    println(&bl,bl)
    println(&cl,cl)
}
```

解答：在go语言中，常量是编译阶段直接使用，不存在地址，取地址有问题。



9.goto跳转的问题

```go
package main

func main()  {

    for i:=0;i<10 ;i++  {
    loop:
        println(i)
    }
    goto loop
}
```

解答：goto 使用有2个限制：1不能跨函数  2.不能跳转到循环的内层代码中。



10. Go 1.9 新特性 Type Alias

    ```go
    package main
    import "fmt"
    
    func main()  {
        type MyInt1 int
        type MyInt2 = int
        var i int = 9
        var i1 MyInt1 = i
        var i2 MyInt2 = i
        fmt.Println(i1,i2)
    }
    ```

    解答：type INT1 = int， INT1 和int 是同一种类型

    ​           type INT   int ，   INT 和 int 不是同一种类型。

```go
package main
import "fmt"

type User struct {
}
type MyUser1 User
type MyUser2 = User
func (i MyUser1) m1(){
    fmt.Println("MyUser1.m1")
}
func (i User) m2(){
    fmt.Println("User.m2")
}

func main() {
    var i1 MyUser1
    var i2 MyUser2
    i1.m1()
    i2.m2()
}
```

输出：MyUser1.m1  User.m2



11.panic 捕捉

```go
package main

import (
    "fmt"
    "reflect"
)

func main1()  {
    defer func() {
       if err:=recover();err!=nil{
           fmt.Println(err)
       }else {
           fmt.Println("fatal")
       }
    }()

    defer func() {
        panic("defer panic")
    }()
    panic("panic")
}

func main()  {
    defer func() {
        if err:=recover();err!=nil{
            fmt.Println("++++")
            f:=err.(func()string)
            fmt.Println(err,f(),reflect.TypeOf(err).Kind().String())
        }else {
            fmt.Println("fatal")
        }
    }()

    defer func() {
        panic(func() string {
            return  "defer panic"
        })
    }()
    panic("panic")
}
```

输出：多个panic， 仅仅会捕捉最后一个。



12.interface

```go
func Foo(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}
func main() {
	var x *int = nil
	Foo(x)
}
```

解答：输出non-empty interface

​           go中interface 有2个部分，type和data， 只有当2个都为nil时， interface才为nil。



13.append 切片操作



```go
package main

import "fmt"

func main() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	s1 = append(s1, s2)
	fmt.Println(s1)
}
```

解答：s1和s2都是切片， 正确的写法是s1=append(s1, s2, ...)



14.interface的类型断言

```go
func main() {
	i := GetValue()

	switch i.(type) {
	case int:
		println("int")
	case string:
		println("string")
	case interface{}:
		println("interface")
	default:
		println("unknown")
	}

}

func GetValue() int {
	return 1
}
```

解答：编译不顾， type不能使用在非interface类型的判断上。将GetValue返回值int改为interface{}即可。

