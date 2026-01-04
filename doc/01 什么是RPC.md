# 什么是RPC

RPC（Remote Procedure Call 远程过程调用）是一种软件通信协议，一个程序可以使用该协议向位于网络上另一台计算机中的程序请求服务，而无需了解网络的详细信息。RPC 用于调用远程系统上的其他进程，如本地系统。过程调用有时也称为 函数调用或 子程序调用。

RPC是一种客户端-服务器交互形式（调用者是客户端，执行者是服务器），通常通过请求-响应消息传递系统实现。与本地过程调用一样，RPC 是一种 同步 操作，需要阻塞请求程序，直到返回远程过程的结果。但是，使用共享相同地址空间的轻量级进程或 线程 可以同时执行多个 RPC。

通俗的解释：客户端在不知道调用细节的情况下，调用存在于远程计算机上的某个对象，就像调用本地应用程序中的对象一样。

接口定义语言（IDL）——用于描述软件组件的应用程序编程接口（API）的规范语言——通常用于远程过程调用软件。在这种情况下，IDL 在链路两端的机器之间提供了一座桥梁，这些机器可能使用不同的操作系统 (OS) 和计算机语言。

实际场景：

有两台服务器，分别是服务器 A、服务器 B。在 服务器 A 上的应用 想要调用服务器 B 上的应用，它们可以直接本地调用吗？

答案是不能的，但走 RPC 的话，十分方便。因此常有人称使用 RPC，就跟本地调用一个函数一样简单。

![019e36d92681455c86542c493fd5af0c.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_e6679925ad674055a14a40a5c2baa996.png "019e36d92681455c86542c493fd5af0c.png")

# HTTP和RPC的区别

RPC关注"方法调用"，HTTP API关注"资源操作"。RPC更像调用函数，HTTP API更像操作数据。

1）概念区别

RPC是一种方法，而HTTP是一种协议。两者都常用于实现服务，在这个层面最本质的区别是RPC服务主要工作在TCP协议之上（也可以在HTTP协议），而HTTP服务工作在HTTP协议之上。由于HTTP协议基于TCP协议，所以RPC服务天然比HTTP更轻量，效率更胜一筹。

两者都是基于网络实现的，从这一点上，都是基于Client/Server架构。

2）从协议上区分

RPC是远端过程调用，其调用协议通常包含：传输协议 和 序列化协议。

* •
  传输协议：著名的 grpc，它底层使用的是 http2 协议；还有 dubbo 一类的自定义报文的 tcp 协议。
* •
  序列化协议：基于文本编码的 json 协议；也有二进制编码的 protobuf、hession 等协议；还有针对 java 高性能、高吞吐量的 kryo 和 ftc 等序列化协议。

HTTP服务工作在HTTP协议之上，而且HTTP协议基于TCP协议。

# RPC如何工作

当调用 RPC 时，调用环境被挂起，过程参数通过网络传送到过程执行的环境，然后在该环境中执行过程。

当过程完成时，结果将被传送回调用环境，在那里继续执行，就像从常规过程调用返回一样。

在 RPC 期间，将执行以下步骤：

1.客户端调用客户端存根。该调用是本地过程调用，参数以正常方式压入堆栈。

2.客户端存根将过程参数打包到消息中并进行系统调用以发送消息。过程参数的打包称为编组。

3.客户端的本地操作系统将消息从客户端机器发送到远程服务器机器。

4.服务器操作系统将传入的数据包传递给服务器存根。

5.服务器存根从消息中解包参数——称为解编组。

6.当服务器过程完成时，它返回到服务器存根，它将返回值编组为一条消息。然后服务器 存根将消息交给传输层。

7.传输层将生成的消息发送回客户端传输层，传输层将消息返回给客户端存根。

8.客户端存根解组返回参数，然后执行返回给调用者。

# RPC的四个核心组件

Client （客户端）：服务调用方。

Server（服务端）：服务提供方。

Client Stub（客户端存根）：存放服务端的地址消息，负责将客户端的请求参数打包成网络消息，然后通过网络发送给服务提供方。

Server Stub（服务端存根）：接收客户端发送的消息，再将客户端请求参数打包成网络消息，然后通过网络远程发送给服务方。

![ab134da94584496d95c0625fd0d8e2d7.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_c94f2c72d542475bbed9b4442de16bec.png "ab134da94584496d95c0625fd0d8e2d7.png")

# RPC的四种调用方式

RPC调用通常根据双端是否流式交互，分为了单项RPC、服务端流式RPC、客户端流式RPC、双向流PRC四种方式。

这里举一个例子，假设你是小超，有一个女朋友叫婷婷，婷婷的每种情绪代表一个微服务，你们之间的每一次对话可以理解为一次PRC调用，为了便于画流程图，RPC请求被封装成client.SayHello，请求包为HelloRequest，响应为HelloReply。

1）单项RPC

即客户端发送一个请求给服务端，从服务端获取一个应答，就像一次普通的函数调用。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
|客户端进程|----->HelloRequest |服务端进程|
|        | <-----HelloReply  |        |
```

* •
  client层调用SayHello接口，把HelloRequest包进行序列化
* •
  client option将序列化的数据发送到server端
* •
  server option接收到RPC请求
* •
  将RPC请求返回给server端，server端进行处理，将结果给server option
* •
  server option将HelloReply进行序列化并发给client option
* •
  client option做反序列化处理，并返回给client层

2）服务端流式RPC

即客户端发送一个请求给服务端，可获取一个数据流用来读取一系列消息。客户端从返回的数据流里一直读取直到没有更多消息为止。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
|客户端进程|----->HelloRequest |服务端进程|
|        | <-----HelloReply  |        |
|        | <-----HelloReply  |        |
|        | <-----...         |        |
|        | <-----RPC函数调用结束|        |
|        | <-----HelloReply  |        |
|        | <-----HelloReply  |        |
```

* •
  client层调用SayHello接口，把HelloRequest包进行序列化
* •
  client option将序列化的数据发送到server端
* •
  server option接收到rpc请求
* •
  将rpc请求返回给server端，server端进行处理，将将数据流给server option
* •
  server option将HelloReply进行序列化并发给
* •
  client client option做反序列化处理，并返回给client层

3）客户端流式RPC

即客户端用提供的一个数据流写入并发送一系列消息给服务端。一旦客户端完成消息写入，就等待服务端读取这些消息并返回应答。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
|客户端进程|----->HelloRequest |服务端进程|
|        |----->HelloRequest |        |
|        |----->...          |        |
|        | <-----SendAndClose|        |
|        | <-----HelloReply  |        |
```

* •
  client层调用SayHello接口，把HelloRequest包进行序列化
* •
  client option将序列化的数据流发送到server端
* •
  server option接收到rpc请求
* •
  将rpc请求返回给server端，server端进行处理，将结果给server option
* •
  server option将HelloReply进行序列化并发给client
* •
  client option做反序列化处理，并返回给client层

4）双向流RPC

双向流 RPC，即两边都可以分别通过一个读写数据流来发送一系列消息。这两个数据流操作是相互独立的，所以客户端和服务端能按其希望的任意顺序读写，例如：服务端可以在写应答前等待所有的客户端消息，或者它可以先读一个消息再写一个消息，或者是读写相结合的其他方式。每个数据流里消息的顺序会被保持。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
|客户端进程|----->HelloRequest |服务端进程|
|        |----->HelloRequest |        |
|        | <-----HelloReply  |        |
|        | <-----HelloReply  |        |
|        | <-----HelloReply  |        |
|        | ----双方都end----  |        |
```

# RPC的优缺点

尽管它拥有广泛的好处，但使用 RPC 的人肯定应该注意一些陷阱。

RPC 为开发人员和应用程序管理人员提供的一些优势：

* •
  帮助客户端通过传统使用高级语言中的过程调用与服务器进行通信。
* •
  可以在分布式环境中使用，也可以在本地环境中使用。
* •
  支持面向进程和面向线程的模型。
* •
  对用户隐藏内部消息传递机制。
* •
  只需极少的努力即可重写和重新开发代码。
* •
  提供抽象，即网络通信的消息传递特性对用户隐藏。
* •
  省略许多协议层以提高性能。

另一方面，RPC 的一些缺点包括：

* •
  客户端和服务器各自的例程使用不同的执行环境，资源（如文件）的使用也更加复杂。因此，RPC 系统并不总是适合传输大量数据。
* •
  RPC 极易发生故障，因为它涉及一个通信系统、另一台机器和另一个进程。
* •
  RPC没有统一的标准；它可以通过多种方式实现。
* •
  RPC 只是基于交互的，因此它在硬件架构方面没有提供任何灵活性。

# 常见的RPC框架

1）跟语言绑定框架

* •
  Dubbo：国内最早开源的 RPC 框架，由阿里巴巴公司开发并于 2011 年末对外开源，仅支持 Java 语言。
* •
  Motan：微博内部使用的 RPC 框架，于 2016 年对外开源，仅支持 Java 语言。
* •
  Tars：腾讯内部使用的 RPC 框架，于 2017 年对外开源，仅支持 C++ 语言。
* •
  Spring Cloud：国外 Pivotal 公司 2014 年对外开源的 RPC 框架，仅支持 Java 语言。

2）跨语言开源框架

* •
  gRPC：Google 于 2015 年对外开源的跨语言 RPC 框架，支持多种语言。
* •
  Thrift：最初是由Facebook 开发的内部系统跨语言的 RPC 框架，2007 年贡献给了 Apache 基金，成为 Apache 开源项目之一，支持多种语言。
* •
  Rpcx：是一个类似阿里巴巴 Dubbo和微博 Motan的 RPC 框架，开源，支持多种语言。

# RPC快速入门

Go语言标准包(net/rpc)已经提供了对RPC的支持，而且支持三个级别的RPC：TCP、HTTP和JSONRPC。但Go语言的RPC包是独一无二的RPC，它和传统的RPC系统不同，它只支持Go语言开发的服务器与客户端之间的交互，因为在内部，它们采用了Gob来编码。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
rpc-tutorial/
├── level1-basic/       # 第1层：基础RPC
├── level2-interface/   # 第2层：接口化RPC
├── level3-jsonrpc/     # 第3层：跨语言JSON-RPC
└── level4-httprpc/     # 第4层：HTTP上的RPC
```

## 1、基础RPC（纯Go，TCP+gob）

1）服务端实现

**level1-basic/server/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "time"
)

// 1. 定义服务结构体
type HelloService struct{}

// 2. 实现RPC方法
// 规则：1. 公开方法 2. 两个参数 3. 第二个是指针 4. 返回error
func (h *HelloService) SayHi(request string, reply *string) error {
    format := time.Now().Format("2006-01-02 15:04:05")
    *reply = fmt.Sprintf("hi %s --- %s", request, format)
    return nil
}

// 3. 计算服务
func (h *HelloService) Add(args [2]int, reply *int) error {
    *reply = args[0] + args[1]
    return nil
}

func main() {
    fmt.Println("🚀 启动基础RPC服务器...")
  
    // 1. 创建服务实例
    helloService := new(HelloService)
  
    // 2. 注册服务
    // 注意：这里用的是 Register，不是 RegisterName
    err := rpc.Register(helloService)
    if err != nil {
        log.Fatal("注册服务失败:", err)
    }
  
    // 3. 监听TCP端口
    listener, err := net.Listen("tcp", ":8888")
    if err != nil {
        log.Fatal("监听失败:", err)
    }
  
    fmt.Println("✅ 服务器启动成功，监听端口 8888")
    fmt.Println("📡 可用服务:")
    fmt.Println("   - HelloService.SayHi(string) -> string")
    fmt.Println("   - HelloService.Add([2]int) -> int")
  
    // 4. 接受连接并处理
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err)
            continue
        }
      
        fmt.Printf("🔗 新连接: %s\n", conn.RemoteAddr())
      
        // 5. 为每个连接启动goroutine处理
        go rpc.ServeConn(conn)
    }
}
```

2）客户端实现

**level1-basic/client/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "fmt"
    "log"
    "net/rpc"
    "time"
)

func main() {
    fmt.Println("🔄 启动基础RPC客户端...")
  
    // 1. 连接RPC服务器
    client, err := rpc.Dial("tcp", "localhost:8888")
    if err != nil {
        log.Fatal("连接服务器失败:", err)
    }
    defer client.Close()
  
    fmt.Println("✅ 成功连接到服务器")
  
    // 2. 测试 SayHi 方法
    for i := 1; i <= 3; i++ {
        var reply string
        err = client.Call("HelloService.SayHi", fmt.Sprintf("用户%d", i), &reply)
        if err != nil {
            log.Printf("调用失败: %v", err)
        } else {
            fmt.Printf("📨 调用 SayHi: %s\n", reply)
        }
        time.Sleep(1 * time.Second)
    }
  
    // 3. 测试 Add 方法
    args := [2]int{100, 200}
    var sum int
    err = client.Call("HelloService.Add", args, &sum)
    if err != nil {
        log.Printf("调用 Add 失败: %v", err)
    } else {
        fmt.Printf("🧮 计算 %d + %d = %d\n", args[0], args[1], sum)
    }
  
    // 4. 测试不存在的服务
    var dummy string
    err = client.Call("NonExistent.Method", "test", &dummy)
    if err != nil {
        fmt.Printf("❌ 预期中的错误（调用不存在的方法）: %v\n", err)
    }
  
    fmt.Println("🎉 客户端测试完成")
}
```

## 2、接口化RPC（设计模式）

在涉及 RPC 的应用中，作为开发人员一般至少有三种角色：首先是服务端实现 RPC 方法的开发人员，其次是客户端调用 RPC 方法的人员，最后也是最重要的是制定服务端和客户端 RPC 接口规范的设计人员。在前面的例子中我们为了简化将以上几种角色的工作全部放到了一起，虽然看似实现简单，但是不利于后期的维护和工作的切割。

1）服务端重构

**level2-interface/shared/rpc\_interface.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package shared

import "net/rpc"

// ===================== 接口规范 =====================
// 这部分由架构师/设计者编写

const (
    // 服务名称（包含包路径，避免冲突）
    HelloServiceName = "tutorial/HelloService"
    MathServiceName  = "tutorial/MathService"
)

// HelloService 接口定义
type HelloServiceInterface interface {
    SayHi(request string, reply *string) error
    Greet(name string, reply *string) error
}

// MathService 接口定义
type MathServiceInterface interface {
    Add(args [2]int, reply *int) error
    Multiply(args [2]int, reply *int) error
}

// ===================== 客户端包装 =====================
// 这部分由客户端开发者编写

// HelloServiceClient 客户端包装
type HelloServiceClient struct {
    *rpc.Client
}

// 创建HelloService客户端
func DialHelloService(network, address string) (*HelloServiceClient, error) {
    client, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &HelloServiceClient{Client: client}, nil
}

// SayHi 客户端方法
func (c *HelloServiceClient) SayHi(request string, reply *string) error {
    return c.Client.Call(HelloServiceName+".SayHi", request, reply)
}

// Greet 客户端方法
func (c *HelloServiceClient) Greet(name string, reply *string) error {
    return c.Client.Call(HelloServiceName+".Greet", name, reply)
}

// MathServiceClient 客户端包装
type MathServiceClient struct {
    *rpc.Client
}

// 创建MathService客户端
func DialMathService(network, address string) (*MathServiceClient, error) {
    client, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &MathServiceClient{Client: client}, nil
}

// Add 客户端方法
func (c *MathServiceClient) Add(a, b int, reply *int) error {
    return c.Client.Call(MathServiceName+".Add", [2]int{a, b}, reply)
}

// Multiply 客户端方法
func (c *MathServiceClient) Multiply(a, b int, reply *int) error {
    return c.Client.Call(MathServiceName+".Multiply", [2]int{a, b}, reply)
}
```

**level2-interface/server/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "time"
  
    "level2-interface/shared"
)

// ===================== 服务实现 =====================
// 这部分由服务端开发者编写

// HelloService 实现
type HelloService struct{}

func (h *HelloService) SayHi(request string, reply *string) error {
    format := time.Now().Format("15:04:05")
    *reply = fmt.Sprintf("[%s] Hi %s!", format, request)
    return nil
}

func (h *HelloService) Greet(name string, reply *string) error {
    hour := time.Now().Hour()
    var greeting string
    switch {
    case hour < 12:
        greeting = "早上好"
    case hour < 18:
        greeting = "下午好"
    default:
        greeting = "晚上好"
    }
    *reply = fmt.Sprintf("%s, %s!", greeting, name)
    return nil
}

// MathService 实现
type MathService struct{}

func (m *MathService) Add(args [2]int, reply *int) error {
    *reply = args[0] + args[1]
    return nil
}

func (m *MathService) Multiply(args [2]int, reply *int) error {
    *reply = args[0] * args[1]
    return nil
}

// ===================== 服务注册 =====================
// 注册服务（使用封装函数）
func registerServices() {
    // 注册 HelloService
    err := rpc.RegisterName(shared.HelloServiceName, new(HelloService))
    if err != nil {
        log.Fatal("注册 HelloService 失败:", err)
    }
  
    // 注册 MathService
    err = rpc.RegisterName(shared.MathServiceName, new(MathService))
    if err != nil {
        log.Fatal("注册 MathService 失败:", err)
    }
}

func main() {
    fmt.Println("🚀 启动接口化RPC服务器...")
  
    // 注册服务
    registerServices()
  
    // 启动TCP监听
    listener, err := net.Listen("tcp", ":8889")
    if err != nil {
        log.Fatal("监听失败:", err)
    }
  
    fmt.Println("✅ 服务器启动成功，监听端口 8889")
    fmt.Println("📡 可用服务:")
    fmt.Printf("   - %s.SayHi(string) -> string\n", shared.HelloServiceName)
    fmt.Printf("   - %s.Greet(string) -> string\n", shared.HelloServiceName)
    fmt.Printf("   - %s.Add([2]int) -> int\n", shared.MathServiceName)
    fmt.Printf("   - %s.Multiply([2]int) -> int\n", shared.MathServiceName)
  
    // 接受连接
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err)
            continue
        }
      
        fmt.Printf("🔗 新连接: %s\n", conn.RemoteAddr())
        go rpc.ServeConn(conn)
    }
}
```

2）客户端重构

**level2-interface/server/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "fmt"
    "log"
    "time"
  
    "level2-interface/shared"
)

func main() {
    fmt.Println("🔄 启动接口化RPC客户端...")
  
    // 1. 连接到 HelloService
    helloClient, err := shared.DialHelloService("tcp", "localhost:8889")
    if err != nil {
        log.Fatal("连接 HelloService 失败:", err)
    }
    defer helloClient.Close()
  
    // 2. 连接到 MathService
    mathClient, err := shared.DialMathService("tcp", "localhost:8889")
    if err != nil {
        log.Fatal("连接 MathService 失败:", err)
    }
    defer mathClient.Close()
  
    fmt.Println("✅ 成功连接到所有服务")
  
    // 3. 测试 HelloService
    fmt.Println("\n🧪 测试 HelloService:")
    for i := 1; i <= 2; i++ {
        var reply string
        err = helloClient.SayHi(fmt.Sprintf("用户%d", i), &reply)
        if err != nil {
            log.Printf("SayHi 失败: %v", err)
        } else {
            fmt.Printf("   SayHi: %s\n", reply)
        }
        time.Sleep(500 * time.Millisecond)
    }
  
    // 测试 Greet
    var greeting string
    err = helloClient.Greet("张三", &greeting)
    if err != nil {
        log.Printf("Greet 失败: %v", err)
    } else {
        fmt.Printf("   Greet: %s\n", greeting)
    }
  
    // 4. 测试 MathService
    fmt.Println("\n🧮 测试 MathService:")
    var sum int
    err = mathClient.Add(25, 75, &sum)
    if err != nil {
        log.Printf("Add 失败: %v", err)
    } else {
        fmt.Printf("   25 + 75 = %d\n", sum)
    }
  
    var product int
    err = mathClient.Multiply(12, 12, &product)
    if err != nil {
        log.Printf("Multiply 失败: %v", err)
    } else {
        fmt.Printf("   12 × 12 = %d\n", product)
    }
  
    fmt.Println("\n🎉 客户端测试完成")
  
    // 5. 演示编译时检查
    // 下面这行代码会编译失败，因为客户端接口限制了我们只能调用定义好的方法
    // helloClient.Call("NonExistent.Method", "test", &dummy) // 编译错误！
}
```

## 3、跨语言JSON-RPC

标准库的RPC默认采用 Go 语言特有的 gob 编码，因此从其他语言调用 Go 语言实现的 RPC 服务将比较困难。在互联网的微服务时代，每个 RPC 以及服务的使用者都可能采用不同的编程语言，因此跨语言是互联网时代 RPC 的一个首要条件。得益于 RPC 的框架设计，Go 语言的 RPC 其实也是很容易实现跨语言支持的。

Go 语言的 RPC 框架有两个比较有特色的设计：

* •
  RPC 数据打包时可以通过插件实现自定义的编码和解码。
* •
  RPC 建立在抽象的 io.ReadWriterCloser 接口之上的，我们可以将 RPC 架设在不同的通信协议之上。

这里我们使用Go官方自带的 net/rpc/jsonrpc 扩展实现一个跨语言的rpc。

1）服务端实现

**level3-jsonrpc/server/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "time"
)

// 用户结构体
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// UserService 用户服务
type UserService struct{}

// GetUserInfo 获取用户信息
func (u *UserService) GetUserInfo(name string, reply *User) error {
    *reply = User{
        Name:  name,
        Email: fmt.Sprintf("%s@example.com", name),
        Age:   25 + len(name)%10,
    }
    return nil
}

// RegisterUser 注册用户
func (u *UserService) RegisterUser(user User, reply *map[string]interface{}) error {
    *reply = map[string]interface{}{
        "status":    "success",
        "message":   "用户注册成功",
        "user_id":   fmt.Sprintf("user_%d", time.Now().Unix()),
        "timestamp": time.Now().Format(time.RFC3339),
        "user_info": user,
    }
    return nil
}

// CalculatorService 计算服务
type CalculatorService struct{}

// Calculate 计算表达式
func (c *CalculatorService) Calculate(expr string, reply *float64) error {
    // 简化的计算逻辑，实际项目中应该用表达式解析器
    switch expr {
    case "1+1":
        *reply = 2
    case "2*3":
        *reply = 6
    case "10/2":
        *reply = 5
    case "3^2":
        *reply = 9
    default:
        *reply = 0
    }
    return nil
}

func main() {
    fmt.Println("🚀 启动JSON-RPC服务器...")
  
    // 注册服务
    err := rpc.RegisterName("UserService", new(UserService))
    if err != nil {
        log.Fatal("注册 UserService 失败:", err)
    }
  
    err = rpc.RegisterName("CalculatorService", new(CalculatorService))
    if err != nil {
        log.Fatal("注册 CalculatorService 失败:", err)
    }
  
    // 启动TCP监听
    listener, err := net.Listen("tcp", ":8890")
    if err != nil {
        log.Fatal("监听失败:", err)
    }
  
    fmt.Println("✅ 服务器启动成功，监听端口 8890")
    fmt.Println("📡 可用服务 (JSON格式):")
    fmt.Println("   UserService.GetUserInfo(string) -> User")
    fmt.Println("   UserService.RegisterUser(User) -> map")
    fmt.Println("   CalculatorService.Calculate(string) -> float64")
    fmt.Println("\n📝 示例请求 (可用 curl 测试):")
    fmt.Println(`   echo '{"method":"UserService.GetUserInfo","params":["Alice"],"id":1}' | nc localhost 8890`)
  
    // 接受连接
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err)
            continue
        }
      
        fmt.Printf("🔗 新连接: %s\n", conn.RemoteAddr())
      
        // 关键：使用 JSON 编解码器
        go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}
```

2）客户端实现

**level3-jsonrpc/client/go\_client.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
)

// 用户结构体（必须和服务端匹配）
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    fmt.Println("🔄 Go JSON-RPC 客户端启动...")
  
    // 1. 建立TCP连接
    conn, err := net.Dial("tcp", "localhost:8890")
    if err != nil {
        log.Fatal("连接失败:", err)
    }
    defer conn.Close()
  
    fmt.Println("✅ 成功连接到服务器")
  
    // 2. 创建JSON-RPC客户端
    client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
    defer client.Close()
  
    // 3. 测试 UserService.GetUserInfo
    fmt.Println("\n🧪 测试 UserService.GetUserInfo:")
    var user User
    err = client.Call("UserService.GetUserInfo", "Alice", &user)
    if err != nil {
        log.Printf("调用失败: %v", err)
    } else {
        userJSON, _ := json.MarshalIndent(user, "   ", "  ")
        fmt.Printf("   用户信息: %s\n", userJSON)
    }
  
    // 4. 测试 UserService.RegisterUser
    fmt.Println("\n📝 测试 UserService.RegisterUser:")
    newUser := User{
        Name:  "Bob",
        Email: "bob@example.com",
        Age:   30,
    }
    var result map[string]interface{}
    err = client.Call("UserService.RegisterUser", newUser, &result)
    if err != nil {
        log.Printf("调用失败: %v", err)
    } else {
        resultJSON, _ := json.MarshalIndent(result, "   ", "  ")
        fmt.Printf("   注册结果: %s\n", string(resultJSON))
    }
  
    // 5. 测试 CalculatorService
    fmt.Println("\n🧮 测试 CalculatorService.Calculate:")
    var answer float64
    expressions := []string{"1+1", "2*3", "10/2", "3^2"}
    for _, expr := range expressions {
        err = client.Call("CalculatorService.Calculate", expr, &answer)
        if err != nil {
            log.Printf("计算 %s 失败: %v", expr, err)
        } else {
            fmt.Printf("   %s = %.0f\n", expr, answer)
        }
    }
  
    fmt.Println("\n🎉 Go 客户端测试完成")
}
```

**level3-jsonrpc/client/python\_client.py**（跨语言演示）

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
#!/usr/bin/env python3
"""
Python JSON-RPC 客户端演示
展示如何用其他语言调用 Go 的 JSON-RPC 服务
"""
import json
import socket
import time

class JSONRPCClient:
    def __init__(self, host='localhost', port=8890):
        self.host = host
        self.port = port
        self.request_id = 1
  
    def call(self, method, params):
        """发送JSON-RPC请求"""
        # 构造请求
        request = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params,
            "id": self.request_id
        }
        self.request_id += 1
      
        # 连接服务器
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((self.host, self.port))
      
        # 发送请求
        request_json = json.dumps(request) + "\n"
        sock.sendall(request_json.encode('utf-8'))
      
        # 接收响应
        response_data = b""
        while True:
            chunk = sock.recv(1024)
            if not chunk:
                break
            response_data += chunk
            if b'\n' in chunk:
                break
      
        sock.close()
      
        # 解析响应
        try:
            response = json.loads(response_data.decode('utf-8').strip())
            if 'error' in response and response['error'] is not None:
                print(f"错误: {response['error']}")
                return None
            return response.get('result')
        except json.JSONDecodeError as e:
            print(f"JSON解析错误: {e}")
            print(f"原始响应: {response_data}")
            return None

def main():
    print("🐍 Python JSON-RPC 客户端启动...")
    client = JSONRPCClient()
  
    # 1. 测试 GetUserInfo
    print("\n1️⃣ 获取用户信息:")
    user = client.call("UserService.GetUserInfo", ["Charlie"])
    if user:
        print(f"   用户: {user['Name']}")
        print(f"   邮箱: {user['Email']}")
        print(f"   年龄: {user['Age']}")
  
    # 2. 测试计算器
    print("\n2️⃣ 测试计算器:")
    expressions = ["1+1", "2*3", "10/2"]
    for expr in expressions:
        result = client.call("CalculatorService.Calculate", [expr])
        if result is not None:
            print(f"   {expr} = {result}")
  
    # 3. 注册新用户
    print("\n3️⃣ 注册新用户:")
    new_user = {
        "Name": "David",
        "Email": "david@example.com",
        "Age": 28
    }
    registration = client.call("UserService.RegisterUser", [new_user])
    if registration:
        print(f"   状态: {registration['status']}")
        print(f"   消息: {registration['message']}")
        print(f"   用户ID: {registration['user_id']}")
  
    print("\n✅ Python 客户端测试完成")

if __name__ == "__main__":
    main()
```

## 4、HTTP 上的 RPC

1）服务端实现

**level4-httprpc/server/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/rpc"
    "net/rpc/jsonrpc"
    "strings"
    "time"
)

// 产品结构体
type Product struct {
    ID     string  `json:"id"`
    Name   string  `json:"name"`
    Price  float64 `json:"price"`
    Stock  int     `json:"stock"`
    Status string  `json:"status"` // "available", "out_of_stock", "discontinued"
}

// 订单结构体
type Order struct {
    OrderID    string          `json:"order_id"`
    UserID     string          `json:"user_id"`
    Products   []OrderItem     `json:"products"`
    Total      float64         `json:"total"`
    Status     string          `json:"status"` // "pending", "paid", "shipped", "delivered"
    CreatedAt  string          `json:"created_at"`
}

type OrderItem struct {
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// ProductService 产品服务
type ProductService struct {
    products map[string]Product
}

func NewProductService() *ProductService {
    return &ProductService{
        products: map[string]Product{
            "P001": {ID: "P001", Name: "笔记本电脑", Price: 5999.99, Stock: 50, Status: "available"},
            "P002": {ID: "P002", Name: "智能手机", Price: 2999.99, Stock: 100, Status: "available"},
            "P003": {ID: "P003", Name: "平板电脑", Price: 1999.99, Stock: 0, Status: "out_of_stock"},
            "P004": {ID: "P004", Name: "智能手表", Price: 999.99, Stock: 30, Status: "available"},
        },
    }
}

// GetProduct 获取产品信息
func (p *ProductService) GetProduct(productID string, reply *Product) error {
    product, exists := p.products[productID]
    if !exists {
        return fmt.Errorf("产品不存在: %s", productID)
    }
    *reply = product
    return nil
}

// ListProducts 列出所有产品
func (p *ProductService) ListProducts(_ string, reply *[]Product) error {
    products := make([]Product, 0, len(p.products))
    for _, product := range p.products {
        products = append(products, product)
    }
    *reply = products
    return nil
}

// OrderService 订单服务
type OrderService struct {
    orders map[string]Order
}

func NewOrderService() *OrderService {
    return &OrderService{
        orders: make(map[string]Order),
    }
}

// CreateOrder 创建订单
func (o *OrderService) CreateOrder(args map[string]interface{}, reply *Order) error {
    userID, _ := args["user_id"].(string)
    items, _ := args["items"].([]interface{})
  
    // 计算总价
    var total float64
    orderItems := make([]OrderItem, 0, len(items))
  
    for i, item := range items {
        itemMap, _ := item.(map[string]interface{})
        productID, _ := itemMap["product_id"].(string)
        quantity, _ := itemMap["quantity"].(float64)
        price, _ := itemMap["price"].(float64)
      
        orderItems = append(orderItems, OrderItem{
            ProductID: productID,
            Quantity:  int(quantity),
            Price:     price,
        })
      
        total += price * quantity
      
        log.Printf("订单项 %d: %s × %d = %.2f", i+1, productID, int(quantity), price*quantity)
    }
  
    orderID := fmt.Sprintf("ORD%08d", time.Now().UnixNano()%100000000)
  
    order := Order{
        OrderID:   orderID,
        UserID:    userID,
        Products:  orderItems,
        Total:     total,
        Status:    "pending",
        CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
    }
  
    o.orders[orderID] = order
    *reply = order
  
    log.Printf("订单创建成功: %s, 总价: %.2f", orderID, total)
    return nil
}

// JSON-RPC 处理器
func jsonrpcHandler(w http.ResponseWriter, r *http.Request) {
    // 设置CORS头
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  
    // 处理预检请求
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
  
    // 只接受POST请求
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
  
    // 设置响应头
    w.Header().Set("Content-Type", "application/json")
  
    // 创建连接适配器
    conn := &struct {
        io.Writer
        io.ReadCloser
    }{
        ReadCloser: r.Body,
        Writer:     w,
    }
  
    // 处理RPC请求
    rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
}

// 健康检查端点
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":    "healthy",
        "service":   "HTTP JSON-RPC Server",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

// API文档端点
func docsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    html := `<!DOCTYPE html>
<html>
<head>
    <title>HTTP JSON-RPC API 文档</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 20px; margin: 10px 0; border-radius: 5px; }
        code { background: #eee; padding: 2px 5px; border-radius: 3px; }
        pre { background: #f8f8f8; padding: 10px; overflow: auto; }
    </style>
</head>
<body>
    <h1>📡 HTTP JSON-RPC API 文档</h1>
  
    <div class="endpoint">
        <h2>📊 健康检查</h2>
        <p><code>GET /health</code></p>
        <pre>curl http://localhost:8891/health</pre>
    </div>
  
    <div class="endpoint">
        <h2>🛒 产品服务</h2>
        <p><strong>获取单个产品</strong></p>
        <pre>curl -X POST http://localhost:8891/jsonrpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "ProductService.GetProduct",
    "params": ["P001"],
    "id": 1
  }'</pre>
      
        <p><strong>列出所有产品</strong></p>
        <pre>curl -X POST http://localhost:8891/jsonrpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "ProductService.ListProducts",
    "params": [""],
    "id": 2
  }'</pre>
    </div>
  
    <div class="endpoint">
        <h2>📦 订单服务</h2>
        <p><strong>创建订单</strong></p>
        <pre>curl -X POST http://localhost:8891/jsonrpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "OrderService.CreateOrder",
    "params": [{
      "user_id": "U1001",
      "items": [
        {"product_id": "P001", "quantity": 1, "price": 5999.99},
        {"product_id": "P002", "quantity": 2, "price": 2999.99}
      ]
    }],
    "id": 3
  }'</pre>
    </div>
  
    <div class="endpoint">
        <h2>🔗 快速测试</h2>
        <button onclick="testGetProduct()">测试获取产品</button>
        <button onclick="testCreateOrder()">测试创建订单</button>
        <div id="result" style="margin-top: 20px;"></div>
    </div>
  
    <script>
        async function testGetProduct() {
            const response = await fetch('/jsonrpc', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    jsonrpc: "2.0",
                    method: "ProductService.GetProduct",
                    params: ["P001"],
                    id: Date.now()
                })
            });
            const data = await response.json();
            document.getElementById('result').innerHTML = 
                '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
        }
      
        async function testCreateOrder() {
            const response = await fetch('/jsonrpc', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    jsonrpc: "2.0",
                    method: "OrderService.CreateOrder",
                    params: [{
                        user_id: "U1001",
                        items: [
                            {product_id: "P001", quantity: 1, price: 5999.99},
                            {product_id: "P004", quantity: 3, price: 999.99}
                        ]
                    }],
                    id: Date.now()
                })
            });
            const data = await response.json();
            document.getElementById('result').innerHTML = 
                '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
        }
    </script>
</body>
</html>`
    fmt.Fprint(w, html)
}

func main() {
    fmt.Println("🚀 启动HTTP JSON-RPC服务器...")
  
    // 创建服务实例
    productService := NewProductService()
    orderService := NewOrderService()
  
    // 注册RPC服务
    err := rpc.RegisterName("ProductService", productService)
    if err != nil {
        log.Fatal("注册 ProductService 失败:", err)
    }
  
    err = rpc.RegisterName("OrderService", orderService)
    if err != nil {
        log.Fatal("注册 OrderService 失败:", err)
    }
  
    // 注册HTTP处理器
    http.HandleFunc("/jsonrpc", jsonrpcHandler)
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/", docsHandler)
  
    // 启动HTTP服务器
    port := ":8891"
    fmt.Printf("✅ 服务器启动成功\n")
    fmt.Printf("📡 监听地址: http://localhost%s\n", port)
    fmt.Printf("🔧 可用端点:\n")
    fmt.Printf("   📄 文档:      http://localhost%s/\n", port)
    fmt.Printf("   ❤️  健康检查:  http://localhost%s/health\n", port)
    fmt.Printf("   ⚡ RPC接口:   http://localhost%s/jsonrpc\n", port)
    fmt.Printf("\n🛒 可用服务:\n")
    fmt.Printf("   - ProductService.GetProduct(productID)\n")
    fmt.Printf("   - ProductService.ListProducts()\n")
    fmt.Printf("   - OrderService.CreateOrder(orderData)\n")
  
    // 启动服务器
    err = http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatal("启动HTTP服务器失败:", err)
    }
}
```

2）客户端实现

**level4-httprpc/client/main.go**

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)

// RPC请求结构体
type JSONRPCRequest struct {
    JSONRPC string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params"`
    ID      int64       `json:"id"`
}

// RPC响应结构体
type JSONRPCResponse struct {
    JSONRPC string          `json:"jsonrpc"`
    Result  json.RawMessage `json:"result,omitempty"`
    Error   *RPCError       `json:"error,omitempty"`
    ID      int64           `json:"id"`
}

// RPC错误结构体
type RPCError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    string `json:"data,omitempty"`
}

// 产品结构体
type Product struct {
    ID     string  `json:"id"`
    Name   string  `json:"name"`
    Price  float64 `json:"price"`
    Stock  int     `json:"stock"`
    Status string  `json:"status"`
}

// 订单项
type OrderItem struct {
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// 订单
type Order struct {
    OrderID   string      `json:"order_id"`
    UserID    string      `json:"user_id"`
    Products  []OrderItem `json:"products"`
    Total     float64     `json:"total"`
    Status    string      `json:"status"`
    CreatedAt string      `json:"created_at"`
}

// HTTP JSON-RPC客户端
type HTTPRPCClient struct {
    BaseURL string
    Client  *http.Client
    ID      int64
}

func NewHTTPRPCClient(baseURL string) *HTTPRPCClient {
    return &HTTPRPCClient{
        BaseURL: baseURL,
        Client: &http.Client{
            Timeout: 10 * time.Second,
        },
        ID: 1,
    }
}

// Call 发送RPC请求
func (c *HTTPRPCClient) Call(method string, params interface{}, result interface{}) error {
    // 构造请求
    request := JSONRPCRequest{
        JSONRPC: "2.0",
        Method:  method,
        Params:  params,
        ID:      c.ID,
    }
    c.ID++
  
    // 序列化请求
    requestBody, err := json.Marshal(request)
    if err != nil {
        return fmt.Errorf("序列化请求失败: %v", err)
    }
  
    // 发送HTTP请求
    resp, err := c.Client.Post(c.BaseURL+"/jsonrpc", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return fmt.Errorf("HTTP请求失败: %v", err)
    }
    defer resp.Body.Close()
  
    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("读取响应失败: %v", err)
    }
  
    // 解析响应
    var rpcResp JSONRPCResponse
    if err := json.Unmarshal(body, &rpcResp); err != nil {
        return fmt.Errorf("解析JSON失败: %v\n原始响应: %s", err, string(body))
    }
  
    // 检查错误
    if rpcResp.Error != nil {
        return fmt.Errorf("RPC错误 [%d]: %s", rpcResp.Error.Code, rpcResp.Error.Message)
    }
  
    // 解析结果
    if result != nil && len(rpcResp.Result) > 0 {
        if err := json.Unmarshal(rpcResp.Result, result); err != nil {
            return fmt.Errorf("解析结果失败: %v\n原始结果: %s", err, string(rpcResp.Result))
        }
    }
  
    return nil
}

func main() {
    fmt.Println("🔄 HTTP JSON-RPC 客户端启动...")
  
    // 创建客户端
    client := NewHTTPRPCClient("http://localhost:8891")
  
    // 1. 测试健康检查
    fmt.Println("\n1️⃣ 测试健康检查:")
    resp, err := http.Get("http://localhost:8891/health")
    if err != nil {
        log.Printf("健康检查失败: %v", err)
    } else {
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        fmt.Printf("   健康状态: %s\n", string(body))
    }
  
    // 2. 获取单个产品
    fmt.Println("\n2️⃣ 获取产品信息:")
    var product Product
    err = client.Call("ProductService.GetProduct", []string{"P001"}, &product)
    if err != nil {
        log.Printf("获取产品失败: %v", err)
    } else {
        fmt.Printf("   产品ID: %s\n", product.ID)
        fmt.Printf("   产品名称: %s\n", product.Name)
        fmt.Printf("   价格: ¥%.2f\n", product.Price)
        fmt.Printf("   库存: %d\n", product.Stock)
        fmt.Printf("   状态: %s\n", product.Status)
    }
  
    // 3. 列出所有产品
    fmt.Println("\n3️⃣ 列出所有产品:")
    var products []Product
    err = client.Call("ProductService.ListProducts", []string{""}, &products)
    if err != nil {
        log.Printf("列出产品失败: %v", err)
    } else {
        fmt.Printf("   共有 %d 个产品:\n", len(products))
        for i, p := range products {
            stockStatus := "✅"
            if p.Stock == 0 {
                stockStatus = "❌"
            }
            fmt.Printf("   %d. %s %s - ¥%.2f (%d件库存)\n", 
                i+1, stockStatus, p.Name, p.Price, p.Stock)
        }
    }
  
    // 4. 创建订单
    fmt.Println("\n4️⃣ 创建订单:")
    orderData := map[string]interface{}{
        "user_id": "U1001",
        "items": []map[string]interface{}{
            {"product_id": "P001", "quantity": 1, "price": 5999.99},
            {"product_id": "P004", "quantity": 2, "price": 999.99},
        },
    }
  
    var order Order
    err = client.Call("OrderService.CreateOrder", []interface{}{orderData}, &order)
    if err != nil {
        log.Printf("创建订单失败: %v", err)
    } else {
        fmt.Printf("   🎉 订单创建成功!\n")
        fmt.Printf("   订单号: %s\n", order.OrderID)
        fmt.Printf("   用户ID: %s\n", order.UserID)
        fmt.Printf("   订单状态: %s\n", order.Status)
        fmt.Printf("   订单总额: ¥%.2f\n", order.Total)
        fmt.Printf("   创建时间: %s\n", order.CreatedAt)
        fmt.Printf("   包含 %d 个商品:\n", len(order.Products))
        for i, item := range order.Products {
            fmt.Printf("     %d. %s × %d = ¥%.2f\n", 
                i+1, item.ProductID, item.Quantity, item.Price*float64(item.Quantity))
        }
    }
  
    // 5. 测试错误情况
    fmt.Println("\n5️⃣ 测试错误处理:")
    var dummy Product
    err = client.Call("ProductService.GetProduct", []string{"NON_EXISTENT"}, &dummy)
    if err != nil {
        fmt.Printf("   ✅ 预期错误: %v\n", err)
    }
  
    fmt.Println("\n🎉 HTTP RPC客户端测试完成")
    fmt.Println("\n💡 提示: 你也可以使用curl测试:")
    fmt.Println(`   curl -X POST http://localhost:8891/jsonrpc \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"ProductService.GetProduct","params":["P002"],"id":100}'`)
}
```

xx
