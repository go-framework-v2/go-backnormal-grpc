# 什么是gRPC？

gRPC 是一个强大的开源 RPC（远程过程调用）框架，用于构建可扩展且快速的 API。它允许客户端和服务器应用程序透明地通信并开发连接的系统。gRPC框架依赖 HTTP/2、协议缓冲区和其他现代技术堆栈来确保最大的 API 安全性、性能和可扩展性。

在 gRPC 中，客户端应用程序可以直接调用不同机器上的服务器应用程序上的方法，就像是本地对象一样，更容易创建分布式应用程序和服务。

与许多 RPC 系统一样，gRPC 基于定义服务的思想，指定可以远程调用的方法及其参数和返回类型。在服务端，服务端实现这个接口并运行一个 gRPC 服务器来处理客户端调用。在客户端，客户端有一个存根（自动生成的文件），它提供与服务器相同的方法。

![65739071138f4cbeb98dec3befa78154.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_b2092122c3504faca5c1115974052743.png "65739071138f4cbeb98dec3befa78154.png")

# gRPC 的历史

2015 年，Google 开发了 gRPC 作为 RPC 框架的扩展，以链接使用不同技术创建的许多微服务。最初，它与 Google 的内部基础设施密切相关，但后来，它被开源并标准化以供社区使用。在其发布的第一年，顶级组织利用它来支持从微服务到 Web、移动和物联网的用例。并在 2017 年因越来越受欢迎而成为云原生计算基金会（CNCF）孵化项目。

# 使用Protobuf

Protobuf 是 Google 的序列化/反序列化协议，可以轻松定义服务和自动生成客户端库。gRPC 使用此协议作为其接口定义语言 (IDL) 和序列化工具集。

* •
  客户端和服务器之间的 gRPC 服务和消息在 proto 文件中定义。
* •
  Protobuf 编译器 protoc 生成客户端和服务器代码，在运行时将 .proto 文件加载到内存中，并使用内存中的模式来序列化/反序列化二进制消息。
* •
  代码生成后，每条消息都会在客户端和远程服务之间进行交换。

为什么使用Protobuf？

使用 Protobuf 进行解析需要更少的 CPU 资源，因为数据被转换为二进制格式，并且编码的消息的大小更轻。因此，消息交换速度更快，即使在 CPU 速度较慢的机器（例如移动设备）中也是如此。

# gRPC架构

在下面的 gRPC 架构图中，我们有 gRPC 客户端和服务器端。在 gRPC 中，每个客户端服务都包含一个存根（自动生成的文件），类似于包含当前远程过程的接口。

gRPC工作流程：

* •
  gRPC 客户端将要发送到服务器的参数对存根进行本地过程调用。
* •
  客户端存根使用 Protobuf 使用编组过程序列化参数，并将请求转发到本地机器中的本地客户端时间库。
* •
  操作系统通过 HTTP/2 协议调用远程服务器机器。
* •
  服务器的操作系统接收数据包并调用服务器存根程序，该程序对接收到的参数进行解码并使用 Protobuf 执行相应的程序调用。
* •
  服务器存根将编码响应发送回客户端传输层。客户端存根取回结果消息并解包返回的参数，然后执行返回给调用者。

![84fd0a66d828452ca6afbae8dd031de5.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_d3bc2367d9d1461e8f36394b7e8f84a4.png "84fd0a66d828452ca6afbae8dd031de5.png")

# 准备工作

使用 Go 来编写 gRPC Server 和 Client，让其互相通讯。在此之上会使用到如下库：

下面示例是在 windows环境中安装。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
google.golang.org/grpc
google.golang.org/protobuf/cmd/protoc-gen-go
google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

1.初始化项目go mod init 项目名称模块管理依赖项

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go mod init go-grpc-examle
```

2.安装protoc：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
➜ go get -u google.golang.org/grpc
```

通过--version命令查看是否安装成功:

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
➜ protoc --version
libprotoc 3.20.1
```

2.使用以下命令为 Go 安装协议编译器插件：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
➜ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
➜ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

验证插件是否安装成功:

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
➜ protoc-gen-go --version
protoc-gen-go.exe v1.28.1

➜ protoc-gen-go-grpc --version
protoc-gen-go-grpc 1.2.0
```

3.更新你的PATH，以便protoc编译器可以找到插件（你不在gopath下创建的项目，这里自行百度改一下PATH）

# 编写gRPC Client and Server

## 1、目录结构

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go-grpc-example
├── client
│   └──hello_client
│       └── client.go
├── proto
│   └──hello
│       └── hello.proto
├── server
│   └──hello_server
│       └── server.go
```

## 2、编写 .proto 文件

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
syntax = "proto3";

// 定义go生成后的包名
option go_package = "./;hello";
package proto;

// 定义入参
message Request {
  string name =1;
}
// 定义返回
message Response {
  string result = 1;
}

// 定义接口
service UserService {
  rpc SayHi(Request) returns (Response);
}
```

## 3、生成Go代码

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 同时生成hello.pb.go 和 hello_grpc.pb.go
➜ protoc --go-grpc_out=. --go_out=. hello.proto
```

当前目录下可以看到生成两个文件：

![f2da544cb2754ce7b4e2c4e923de1d9c.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_bc55946caad94237b644b5ef9f5368bb.png "f2da544cb2754ce7b4e2c4e923de1d9c.png")

## 4、编写 Server 服务端代码

编写 gRPC Server 的基础模板，完成一个方法的调用。对 server.go 写入如下内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
	"context"
	"fmt"
	"go-grpc-example/proto/hello"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type HelloService struct {
	// 必须嵌入UnimplementedUserServiceServer
	hello.UnimplementedUserServiceServer
}

// 实现SayHi方法
func (h *HelloService) SayHi(ctx context.Context, req *hello.Request) (res *hello.Response, err error) {
	format := time.Now().Format("2006-01-02 15:04:05")
	return &hello.Response{Result: "hi " + req.GetName() + "---" + format}, nil
}

const PORT = "8888"

func main() {
	PORT := "8888"

	fmt.Println("🚀 启动 gRPC 服务器...")
	fmt.Printf("📡 监听端口: %s\n", PORT)

	// 创建grpc服务
	server := grpc.NewServer()

	// 注册服务
	hello.RegisterUserServiceServer(server, &HelloService{})

	// 监听端口
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("❌ 监听失败: %v", err)
	}

	fmt.Println("✅ gRPC 服务器启动成功")
	fmt.Println("等待客户端连接...")

	// 启动优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		fmt.Printf("\n📢 收到信号: %v，正在关闭服务...\n", sig)
		server.GracefulStop()
		fmt.Println("👋 服务已关闭")
	}()

	// 启动服务
	if err := server.Serve(lis); err != nil {
		log.Fatalf("❌ 服务启动失败: %v", err)
	}
}
```

## 5、编写Client客户端代码

接下来编写 gRPC Go Client 的基础模板，打开 hello\_client/client.go 文件，写入以下内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go-grpc-example/proto/hello"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    fmt.Println("🔄 启动 gRPC 客户端...")
  
    // 连接服务器
    conn, err := grpc.Dial("localhost:8888", 
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock())
    if err != nil {
        log.Fatalf("❌ 连接失败: %v", err)
    }
    defer conn.Close()
  
    // 创建客户端
    client := hello.NewUserServiceClient(conn)
  
    fmt.Println("✅ 连接成功，开始测试...")
  
    // 测试3次调用
    for i := 1; i <= 3; i++ {
        // 创建请求
        name := fmt.Sprintf("用户%d", i)
        req := &hello.Request{Name: name}
      
        // 设置超时
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
      
        // 调用服务
        fmt.Printf("📨 发送请求: name=%s\n", name)
        resp, err := client.SayHi(ctx, req)
        if err != nil {
            log.Printf("❌ 调用失败: %v", err)
            continue
        }
      
        fmt.Printf("📬 收到响应: %s\n", resp.Result)
      
        // 等待1秒
        time.Sleep(1 * time.Second)
    }
  
    fmt.Println("🎉 客户端测试完成")
}
```

## 6、启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:50970

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:51040
resp: hi lin钟一---2022-11-01 14:54:01
```
