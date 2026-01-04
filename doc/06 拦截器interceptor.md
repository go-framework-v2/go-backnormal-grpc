# 一、拦截器概述

## 1、什么是拦截器？

在常规的 HTTP 服务器中，我们可以设置有一个中间件将我们的处理程序包装在服务器上。此中间件可用于在实际提供正确内容之前执行服务器想要执行的任何操作，它可以是身份验证或日志记录或任何东西。

中间件：

中间件供系统软件和应用软件之间连接、便于软件各部件之间的沟通的计算机软件，相当于不同技术、工具和数据库之间的桥梁，例如他可以记录响应时长、记录请求和响应数据日志，身份验证等。

中间件可以在拦截到发送给 handler 的请求，且可以拦截 handler 返回给客户端的响应

gRPC 不同，它允许在服务器和客户端都使用拦截器。

* •
  服务器端拦截器是 gRPC 服务器在到达实际 RPC 方法之前调用的函数。它可以用于多种用途，例如日志记录、跟踪、速率限制、身份验证和授权。
* •
  同样，客户端拦截器是 gRPC 客户端在调用实际 RPC 之前调用的函数。

![76e779120cf149cf8e2c53d37cf28d0e.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_47df594415ee4fbdaf31e7d048589ba2.png "76e779120cf149cf8e2c53d37cf28d0e.png")

## 2、gRPC 拦截器核心概念

* •
  一元是我们大多数人使用的。就是发送一个请求并获得一个响应。
* •
  流是当您发送或接收 protobuf 消息的数据管道时。这意味着如果一个 gRPC 服务响应一个流，消费者可以在这个流中期望多个响应。

具体可以看我讲[《gRPC流》](https://blog.csdn.net/weixin_46618592/article/details/127639689?spm=1001.2014.3001.5502) 其中提到了流和一元的详细概念

拦截器正如他名字的含义，它在 API 请求被执行之前拦截它们。这可用于记录、验证或在处理 API 请求之前发生的任何事情，拦截器还可以做统一接口的认证工作，不需要每一个接口都做一次认证了，多个接口多次访问，只需要在统一个地方认证即可。使用 HTTP API，这在 Golang 中很容易，你可以使用中间件包装 HTTP 处理程序。

gRPC有两种数据通信方式，那必然有两种拦截器：

* •
  UnaryInterceptors — 用于 API 调用，即一个客户端请求和一个服务器响应。
* •
  StreamInterceptors —用于 API 调用，其中客户端发送请求但接收回数据流，允许服务器随时间响应多个项目。实际上，由于 gRPC 是双向的，因此客户端也可以使用它来发送数据。

## 3、服务端拦截器和客户端拦截器

gRPC允许在客户端和服务器以及一元和流式调用中使用拦截器，上面提到过gRPC允许在服务器和客户端都使用拦截器，那么我们就有 4 种不同的拦截器。

如果我们去[go-grpc](https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc)库看看他们是如何处理这个的，我们可以看到四个不同的用例。两种拦截器类型都可用于服务器和客户端。

* •
  UnaryClientInterceptor — 在客户端拦截所有一元 gRPC 调用。
* •
  UnaryServerInterceptor — 在服务器端拦截一元 gRPC 调用。
* •
  StreamClientInterceptor — 拦截器在创建客户端流时触发。
* •
  StreamServerInterceptor — 拦截器在服务器上执行 Stream 之前触发。

![5364b24f79294aafbf5c248e7e647f6d.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_2b549ff5399a4d35845afd78f9287289.png "5364b24f79294aafbf5c248e7e647f6d.png")

关于gRPC拦截器类型的定义：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
type UnaryClientInterceptor func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error

type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)

type StreamClientInterceptor func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)

type StreamServerInterceptor func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
```

## 4、Metadata 元数据

gRPC 允许发送自定义元数据。元数据是键值的一个非常简单的概念。

如果我们查看[golang 元数据规范](https://pkg.go.dev/google.golang.org/grpc/metadata#MD)，我们可以看到它是一个map[string][]string。

元数据可以作为header或trailer发送

* •
  header应该在数据之前发送。
* •
  trailer应在处理完毕后发送。

元数据允许我们在不更改 protobuf 消息的情况下向请求中添加数据。这通常用于添加与请求相关但不属于请求的数据。

例如，我们可以在请求的元数据中添加 JWT 令牌作为身份验证。这允许我们在不改变实际服务器逻辑的情况下使用逻辑扩展 API 端点。这对于身份验证、速率限制或日志记录很有用。

理论够了！我相信我们已经准备好开始测试它了。

# 二、拦截器的使用

## 1、目录结构

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go-grpc-example
├─client
│  ├─hello_client
│  │   └──client.go
│  ├─stream_client
│  │   └──client.go
├─pkg
│  ├─Interceptor
│  │   └──Interceptor.go
├─proto
│  ├─hello
│  ├─stream
└─server
    ├─hello_server
    │  └──server.go
    ├─stream_server
    │  └──server.go
```

偷个懒，这里我们就拿之前的一元和流的示例上使用拦截器

示例我们做一个简单的 interceptor 示例，显示拦截器调用RPC方法前的时间、当前运行程序的操作系统、RPC方法结束后的时间，以及调用RPC的方法名。

创建pkg/Interceptor目录，在Interceptor.go文件里我们写拦截器的方法。

## 2、一元拦截器

### 1）UnaryClientInterceptor

作用：这是我们可以使用客户端元数据丰富消息的地方，例如有关客户端运行的硬件或操作系统的一些信息，或者可能启动我们的跟踪流程。

客户端一元拦截器类型为 grpc.UnaryClientInterceptor，具体如下

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{},
        cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        // 预处理(pre-processing)
        start := time.Now()
        // 获取正在运行程序的操作系统
        cos := runtime.GOOS
        // 将操作系统信息附加到传出请求
        ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)

        // 可以看做是当前 RPC 方法，一般在拦截器中调用 invoker 能达到调用 RPC 方法的效果，当然底层也是 gRPC 在处理。
        // 调用RPC方法(invoking RPC method)
        err := invoker(ctx, method, req, reply, cc, opts...)

        // 后处理(post-processing)
        end := time.Now()
        log.Printf("RPC: %s,,client-OS: '%v' req:%v start time: %s, end time: %s, err: %v", method, cos, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
        return err
    }
}
```

invoker(ctx, method, req, reply, cc, opts...) 是真正调用 RPC 方法。因此我们可以在调用前后增加自己的逻辑：比如调用前检查一下参数之类的，调用后记录一下本次请求处理耗时等。

所谓的拦截器其实就是一个函数，可以分为预处理(pre-processing)、调用RPC方法(invoking RPC method)、后处理(post-processing)三个阶段。

* •
  ctx：Go语言中的上下文，一般和 Goroutine 配合使用，起到超时控制的效果
* •
  method：当前调用的 RPC 方法名
* •
  req：本次请求的参数，只有在处理前阶段修改才有效
* •
  reply：本次请求响应，需要在处理后阶段才能获取到
* •
  cc：gRPC 连接信息
* •
  invoker：可以看做是当前 RPC 方法，一般在拦截器中调用 invoker 能达到调用 RPC 方法的效果，当然底层也是 gRPC 在处理。
* •
  opts：本次调用指定的 options 信息

作为一个客户端拦截器，可以在处理前检查 req 看看本次请求带没带 token 之类的鉴权数据，没有的话就可以在拦截器中加上。

**hello\_client**

建立连接时通过 grpc.WithUnaryInterceptor 指定要加载的拦截器：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
//添加一元拦截器
conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(),
        grpc.WithUnaryInterceptor(Interceptor.UnaryClientInterceptor()))
```

### 2）UnaryServerInterceptor

作用：我们可能想要对请求的真实性进行一些检查，例如对其进行授权，或者检查某些字段是否存在/验证请求。

客户端拦截器与服务端拦截器类似：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (resp interface{}, err error) {
        // 预处理(pre-processing)
        start := time.Now()
        // 从传入上下文获取元数据
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, fmt.Errorf("couldn't parse incoming context metadata")
        }

        // 检索客户端操作系统，如果它不存在，则此值为空
        os := md.Get("client-os")
        // 获取客户端IP地址
        ip, err := getClientIP(ctx)
        if err != nil {
            return nil, err
        }

        // RPC 方法真正执行的逻辑
        // 调用RPC方法(invoking RPC method)
        m, err := handler(ctx, req)
        end := time.Now()
        // 记录请求参数 耗时 错误信息等数据
        // 后处理(post-processing)
        log.Printf("RPC: %s,client-OS: '%v' and IP: '%v' req:%v start time: %s, end time: %s, err: %v", info.FullMethod, os, ip, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
        return m, err
    }
}

// GetClientIP检查上下文以检索客户机的ip地址
func getClientIP(ctx context.Context) (string, error) {
    p, ok := peer.FromContext(ctx)
    if !ok {
        return "", fmt.Errorf("couldn't parse client IP address")
    }
    return p.Addr.String(), nil
}
```

handler(ctx, req) 是真正执行 RPC 方法，与invoker的调用不一样，不要搞混了。因此我们可以在真正执行前后检查数据：比如查看客户端操作系统和客户端IP地址、记录请求参数，耗时，错误信息等数据。

参数具体含义如下：

* •
  ctx：请求上下文
* •
  req：RPC 方法的请求参数
* •
  info：RPC 方法的所有信息
* •
  handler：RPC 方法真正执行的逻辑

**hello\_server**

服务端则是在 NewServer 时指定拦截器：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
//添加一元拦截器
server := grpc.NewServer(grpc.UnaryInterceptor(Interceptor.UnaryServerInterceptor()))
```

### 3）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:51081
2022/11/07 18:37:28 RPC: /hello.UserService/SayHi,client-OS: '[windows]' and IP: '127.0.0.1:51104' req:name:"lin钟一" start time: 2022-11-07T18:37:28+08:00
, end time: 2022-11-07T18:37:28+08:00, err: <nil>

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:51102
2022/11/07 18:37:28 RPC: /hello.UserService/SayHi,,client-OS: 'windows' req:name:"lin钟一" start time: 2022-11-07T18:37:28+08:00, end time: 2022-11-07T18:3
7:28+08:00, err: <nil>
resp: hi lin钟一---2022-11-07 18:37:28
```

## 3、流式拦截器

流拦截器过程和一元拦截器有所不同，同样可以分为3个阶段：

* •
  1）预处理(pre-processing)
* •
  2）调用RPC方法(invoking RPC method)
* •
  3）后处理(post-processing)

预处理阶段和一元拦截器类似，但是调用RPC方法和后处理这两个阶段则完全不同。

StreamAPI 的请求和响应都是通过 Stream 进行传递的，更进一步是通过 Streamer 调用 SendMsg 和 RecvMsg 这两个方法获取的。

然后 Streamer 又是调用RPC方法来获取的，所以在流拦截器中我们可以对 Streamer 进行包装，然后实现 SendMsg 和 RecvMsg 这两个方法。

### 1）StreamClientInterceptor

作用：例如，如果我们将 100 个对象的列表传输到服务器，例如文件或视频的块，我们可以在发送每个块之前拦截，并验证校验和等内容是否有效，将元数据添加到帧等。

本例中通过结构体嵌入的方式，对 Streamer 进行包装，在 SendMsg 和 RecvMsg 之前打印出具体的值。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func StreamClientInterceptor() grpc.StreamClientInterceptor {
    return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
        method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
        log.Printf("opening client streaming to the server method: %v", method)
        // 调用Streamer函数，获得ClientStream
        stream, err := streamer(ctx, desc, cc, method)
        return newStreamClient(stream), err
    }
}

// 嵌入式 streamClient 允许我们访问SendMsg和RecvMsg函数
type streamClient struct {
    grpc.ClientStream
}

// 对ClientStream进行包装
func newStreamClient(c grpc.ClientStream) grpc.ClientStream {
    return &streamClient{c}
}

// RecvMsg从流中接收消息
func (e *streamClient) RecvMsg(m interface{}) error {
    // 在这里，我们可以对接收到的消息执行额外的逻辑，例如
    // 验证
    log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
    if err := e.ClientStream.RecvMsg(m); err != nil {
        return err
    }
    return nil
}

// RecvMsg从流中接收消息
func (e *streamClient) SendMsg(m interface{}) error {
    // 在这里，我们可以对接收到的消息执行额外的逻辑，例如
    // 验证
    log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
    if err := e.ClientStream.SendMsg(m); err != nil {
        return err
    }
    return nil
}
```

因为SendMsg 和 RecvMsg 方法 ClientStream接口内的方法，我们需要先调用 streamer(ctx, desc, cc, method)函数获取到ClientStream再对他进一步结构体封装，实现他SendMsg 和 RecvMsg 方法。

**stream\_client**

通过 grpc.WithStreamInterceptor 指定要加载的拦截器

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(), grpc.WithStreamInterceptor(Interceptor.StreamClientInterceptor()))
if err != nil {
    log.Fatalf("grpc.Dial err: %v", err)
}
defer conn.Close()
```

### 2） StreamServerInterceptor

作用：例如，如果我们正在接收上述文件块，也许我们想确定在传输过程中没有丢失任何内容，并在存储之前再次验证校验和。

与客户端类似：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func StreamServerInterceptor() grpc.StreamServerInterceptor {
    return func(srv interface{}, ss grpc.ServerStream,
        info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
        wrapper := newStreamServer(ss)
        return handler(srv, wrapper)
    }
}

// 嵌入式EdgeServerStream允许我们访问RecvMsg函数
type streamServer struct {
    grpc.ServerStream
}

func newStreamServer(s grpc.ServerStream) grpc.ServerStream {
    return &streamServer{s}
}

// RecvMsg从流中接收消息
func (e *streamServer) RecvMsg(m interface{}) error {
    // 在这里，我们可以对接收到的消息执行额外的逻辑，例如
    // 验证
    log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
    if err := e.ServerStream.RecvMsg(m); err != nil {
        return err
    }
    return nil
}

// RecvMsg从流中接收消息
func (e *streamServer) SendMsg(m interface{}) error {
    // 在这里，我们可以对接收到的消息执行额外的逻辑，例如
    // 验证
    log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
    if err := e.ServerStream.SendMsg(m); err != nil {
        return err
    }
    return nil
}
```

StreamServerInterceptor 拦截器自带 ServerStream 参数，我们直接同样的形式进行结构体嵌入封装，在实现他的方法。

**stream\_server**

通过 grpc.StreamInterceptor 指定要加载的拦截器

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
server := grpc.NewServer(grpc.StreamInterceptor(Interceptor.StreamServerInterceptor()))
```

### 3）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:54096
2022/11/07 19:07:17 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:17+08:00
2022/11/07 19:07:17 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111            
2022/11/07 19:07:17 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:17+08:00
2022/11/07 19:07:17 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:17+08:00  
2022/11/07 19:07:18 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/07 19:07:18 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:18+08:00
2022/11/07 19:07:18 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:18+08:00  
2022/11/07 19:07:19 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/07 19:07:19 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:19+08:00
2022/11/07 19:07:19 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:19+08:00  
2022/11/07 19:07:20 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/07 19:07:20 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:20+08:00
2022/11/07 19:07:20 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:20+08:00
2022/11/07 19:07:21 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/07 19:07:22 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:22+08:00
2022/11/07 19:07:22 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:22+08:00
2022/11/07 19:07:23 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/07 19:07:23 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:23+08:00
2022/11/07 19:07:23 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:23+08:00
2022/11/07 19:07:24 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/07 19:07:24 Receive a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:24+08:00
2022/11/07 19:07:24 Send a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:24+08:00

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:54108
2022/11/07 19:07:17 opening client streaming to the server method: /proto.StreamService/Route
2022/11/07 19:07:17 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:17+08:00  
2022/11/07 19:07:17 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:17+08:00
2022/11/07 19:07:17 resp: pj.name: gRPC Stream Server: Route, pt.value: 0
2022/11/07 19:07:17 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:17+08:00
2022/11/07 19:07:18 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:18+08:00
2022/11/07 19:07:18 resp: pj.name: gRPC Stream Server: Route, pt.value: 1
2022/11/07 19:07:18 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:18+08:00
2022/11/07 19:07:19 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:19+08:00
2022/11/07 19:07:19 resp: pj.name: gRPC Stream Server: Route, pt.value: 2
2022/11/07 19:07:19 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:19+08:00
2022/11/07 19:07:20 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:20+08:00
2022/11/07 19:07:20 resp: pj.name: gRPC Stream Server: Route, pt.value: 3
2022/11/07 19:07:20 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:20+08:00
2022/11/07 19:07:21 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:21+08:00
2022/11/07 19:07:22 resp: pj.name: gRPC Stream Server: Route, pt.value: 4
2022/11/07 19:07:22 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:22+08:00
2022/11/07 19:07:23 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:23+08:00
2022/11/07 19:07:23 resp: pj.name: gRPC Stream Server: Route, pt.value: 5
2022/11/07 19:07:23 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:23+08:00
2022/11/07 19:07:24 Send a message (Type: *stream.StreamRequest) at 2022-11-07T19:07:24+08:00
2022/11/07 19:07:24 resp: pj.name: gRPC Stream Server: Route, pt.value: 6
2022/11/07 19:07:24 Receive a message (Type: *stream.StreamResponse) at 2022-11-07T19:07:24+08:00
Server Closed
```

## 4、实现多个拦截器

gRPC框架中只能为每个服务一起配置一元和流拦截器，，gRPC 会根据不同方法选择对应类型的拦截器执行，因此所有的工作只能在一个函数中完成。

开源的grpc-ecosystem项目中的[go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)包已经基于gRPC对拦截器实现了链式拦截的支持。

### 1）Interceptor 新增一个一元客户端拦截器：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func UnaryClientInterceptorTwo() grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{},
        cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        fmt.Println("我是第二个拦截器")
        // 可以看做是当前 RPC 方法，一般在拦截器中调用 invoker 能达到调用 RPC 方法的效果，当然底层也是 gRPC 在处理。
        // 调用RPC方法(invoking RPC method)
        _ = invoker(ctx, method, req, reply, cc, opts...)
        return nil
    }
}
```

### 2）Client 使用go-grpc-middleware实现链式拦截器：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(),
        grpc.WithUnaryInterceptor(
            // 按照顺序依次执行截取器
            grpc_middleware.ChainUnaryClient(Interceptor.UnaryClientInterceptor(),
                Interceptor.UnaryClientInterceptorTwo()),
        ))
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }
    defer conn.Close()
```

### 3）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:55823
2022/11/07 19:32:28 RPC: /hello.UserService/SayHi,client-OS: '[windows]' and IP: '127.0.0.1:55829' req:name:"lin钟一" start time: 2022-11-07T19:32:28+08:00
, end time: 2022-11-07T19:32:28+08:00, err: <nil>

# 启动客户端
$ go run client.go 
我是第一个拦截器
我是第二个拦截器
2022/11/07 19:32:28 RPC: /hello.UserService/SayHi,,client-OS: 'windows' req:name:"lin钟一" start time: 2022-11-07T19:32:28+08:00, end time: 2022-11-07T19:3
2:28+08:00, err: <nil>
resp: hi lin钟一---2022-11-07 19:32:28
```

# 三、小结

1、拦截器分类与定义 gRPC 拦截器可以分为：一元拦截器和流拦截器，服务端拦截器和客户端拦截器。一共有以下4种类型:

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
grpc.UnaryServerInterceptor
grpc.StreamServerInterceptor
grpc.UnaryClientInterceptor
grpc.StreamClientInterceptor
```

拦截器本质上就是一个特定类型的函数，所以实现拦截器只需要实现对应类型方法（方法签名相同）即可。

2、拦截器执行过程

一元拦截器

* •
  1）预处理
* •
  2）调用RPC方法
* •
  3）后处理

流拦截器

* •
  1）预处理
* •
  2）调用RPC方法 获取 Streamer
* •
  3）后处理

。调用 SendMsg 、RecvMsg 之前

。调用 SendMsg 、RecvMsg

。调用 SendMsg 、RecvMsg 之后

3、拦截器使用及执行顺序

配置多个拦截器时，会按照参数传入顺序依次执行

所以，如果想配置一个 Recovery 拦截器则必须放在第一个，放在最后则无法捕获前面执行的拦截器中触发的 panic。
