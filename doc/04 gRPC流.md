# gRPC的请求模型

gRPC 有两种类型的请求模型：

* •
  一元 - 直接的请求响应映射在 HTTP/2 请求响应之上。

简单来说一元就是一个简单的 RPC，其中客户端使用存根向服务器发送请求并等待响应返回，就像正常的函数调用一样。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
rpc SayHi(Request) returns (Response);
```

流式传输——多个请求和响应通过长寿命 HTTP/2 流进行交换，可以是单向或双向的。

![7f623eeb2d6046ac90cc0cfb77acea40.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_316bc9753cb74c97abdd4ec70bc5031b.png "7f623eeb2d6046ac90cc0cfb77acea40.png")

* •
  其中许多进程可以通过 HTTP/2 的多路复用能力（通过单个 TCP 连接一起发送多个响应或接收多个请求）在单个请求中发生。
* •
  Server-side streaming RPC—— 客户端向服务器发送单个请求并接收回数据序列流（读回一系列消息）。客户端从返回的流中读取，直到没有更多消息为止。
* •
  Client-side streaming RPC—— 客户端向服务器发送数据序列流（写入一系列消息），一旦客户端完成了消息的写入，它会等待服务器读取所有消息并返回其响应结果。
* •
  Bidirectional streaming RPC—— 它是双向流式传输，客户端和服务器使用读写流发送一系列消息。两个流独立运行；因此，因此客户端和服务器可以按照他们喜欢的任何顺序读取和写入。保留每个流中消息的顺序。例如，服务器可以在写入响应之前等待接收所有客户端消息，或者它可以交替读取消息然后写入消息，或其他一些读取和写入的组合。

Server-side streaming RPC：服务器端流式 RPC

Client-side streaming RPC：客户端流式 RPC

Bidirectional streaming RPC：双向流式 RPC

stream可以通过将关键字放在请求类型之前来指定流式处理方法。

# HTTP/2

gRPC 是基于HTTP/2开发的，该协议于 2015 年发布，以克服 HTTP/1.1 的限制。在兼容 HTTP/1.1 的同时，我们来了解一下HTTP/2 带来了许多高级功能，例如：

* •
  二进制分帧层 —— 与 HTTP/1.1 不同，HTTP/2 请求/响应分为小消息并以二进制格式分帧，使消息传输高效。通过二进制帧，HTTP/2 协议使请求/响应多路复用成为可能，而不会阻塞网络资源。
* •
  流式传输 —— 客户端可以请求并且服务器可以同时响应的双向全双工流式传输。
* •
  流控制 —— HTTP/2 中使用流控制机制，可以对用于缓冲动态消息的内存进行详细控制。
* •
  标头压缩 —— HTTP/2 中的所有内容，包括标头，都在发送前进行编码，显着提高了整体性能。使用 HPACK 压缩方式，HTTP/2 只共享与之前的 HTTP 头包不同的值。
* •
  处理 —— 使用 HTTP/2，gRPC 支持同步和异步处理，可用于执行不同类型的交互和流式 RPC。

![2cd2ec33ea6142e08eb1117371c97986.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_439d4bcfd2f34db7b61b822ebc2f1b97.png "2cd2ec33ea6142e08eb1117371c97986.png")

HTTP/2 的所有这些特性使 gRPC 能够使用更少的资源，从而减少在云中运行的应用程序和服务之间的响应时间，并延长运行移动设备的客户端的电池寿命。

# gRPC Streaming, Client and Server

## 1、为什么我们要用流式传输，简单的一元RPC不行么？

流式为什么要存在呢？我们在使用一元请求的时候可能会遇到以下问题：

* •
  数据包过大会造成的瞬时压力。
* •
  接收数据包时，需要所有数据包都接受成功且正确后，才能够回调响应，进行业务处理（无法客户端边发送，服务端边处理）

而流式传输却可以：

* •
  HTTP2 通过长期 TCP 连接多路复用流，因此新请求没有 TCP 连接开销。HTTP2 成帧允许在单个 TCP 数据包中发送多个 gRPC 消息。
* •
  对于长期连接，流式请求应该在每条消息的基础上具有最佳性能。一元请求需要为每个请求建立一个新的 HTTP2 流，包括通过网络发送的附加标头帧。一旦建立，通过流式请求发送的每条新消息只需要通过连接发送消息的数据帧。

## 2、目录结构

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go-grpc-example
├── client
│   └──hello_client
│   │   └── client.go
│   └── stream_client
│       └── client.go
├── proto
│   └──hello
│   │   └── hello.proto
│   └──stream
│   │   └── stream.proto
├── server
│   └──hello_server
│   │   └── server.go
│   └──stream_server
│   │   └── server.go
├── Makefile
```

增加 stream\_server、stream\_client 存放服务端和客户端文件，proto/stream/stream.proto 用于编写 IDL

## 3、编写IDL

在 proto/stream 文件夹下的 stream.proto 文件中，写入如下内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
syntax = "proto3";

option go_package="./;stream";
package proto;

service StreamService {
  //List：服务器端流式 RPC
  rpc List(StreamRequest) returns (stream StreamResponse) {};
  //Record：客户端流式 RPC
  rpc Record(stream StreamRequest) returns (StreamResponse) {};
  //Route：双向流式 RPC
  rpc Route(stream StreamRequest) returns (stream StreamResponse) {};
}


message StreamPoint {
  string name = 1;
  int32 value = 2;
}

message StreamRequest {
  StreamPoint pt = 1;
}

message StreamResponse {
  StreamPoint pt = 1;
}
```

注意关键字 stream，声明其为一个流方法。这里共涉及三个方法，对应关系为

* •
  List：服务器端流式 RPC
* •
  Record：客户端流式 RPC
* •
  Route：双向流式 RPC

## 4、Makefile

这是我拖了很久的关于Makefile的用法，感觉Makefile更适合在项目使用中穿插讲解一下。

有一篇很不错的Makefile文档：[点击跳转](https://makefiletutorial.com/#makefile-cookbook)

作用：Makefile 用于帮助决定大型程序的哪些部分需要重新编译。

这里我们用make gen指令代替proto插件从我们的.proto 服务定义中生成 gRPC 客户端和服务器接口。

在Makefile文件中写入：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
gen:
    protoc --go_out=. --go-grpc_out=. ./proto/stream/*.proto
```

用make gen指令生成Go代码：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
➜ make gen
protoc --go_out=. --go-grpc_out=. ./proto/stream/*.proto
```

![d2a4ceded6164a14967cee532c28c7bc.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_5018b2a34b05421980c26f367fcf0497.png "d2a4ceded6164a14967cee532c28c7bc.png")

注意使用Makefile生成的时候，要注意.proto文件 go\_package 指定生成的位置。

## 5、写出基础模板和空定义

我们先把基础的模板和空定义写出来在进行完善，不太懂的看我上一篇文章

### 1）server.go

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
type StreamService struct {
    pb.UnimplementedStreamServiceServer
}

const PORT = "8888"

func main() {
    server := grpc.NewServer() //创建 gRPC Server 对象
    pb.RegisterStreamServiceServer(server, &StreamService{})

    lis, err := net.Listen("tcp", ":"+PORT)
    if err != nil {
        log.Fatalf("net.Listen err: %v", err)
    }

    server.Serve(lis)
}

//服务端流式RPC，Server是Stream，Client为普通RPC请求
//客户端发送一次普通的RPC请求，服务端通过流式响应多次发送数据集
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
    return nil
}

//客户端流式RPC，单向流
//客户端通过流式多次发送RPC请求给服务端，服务端发送一次普通的RPC请求给客户端
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
    return nil
}

//双向流，由客户端发起流式的RPC方法请求，服务端以同样的流式RPC方法响应请求
//首个请求一定是client发起，具体交互方法（谁先谁后，一次发多少，响应多少，什么时候关闭）根据程序编写方式来确定（可以结合协程）
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
    return nil
}
```

### 2）client.go

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
const PORT = "8888"

func main() {
    conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }
    defer conn.Close()

    client := pb.NewStreamServiceClient(conn)

    err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 1234}})
    if err != nil {
        log.Fatalf("printLists.err: %v", err)
    }

    err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 9999}})
    if err != nil {
        log.Fatalf("printRecord.err: %v", err)
    }

    err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 1111}})
    if err != nil {
        log.Fatalf("printRoute.err: %v", err)
    }
}

func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    return nil
}
```

## 6、Server-side streaming RPC：服务器端流式 RPC

服务端流式RPC，Server是Stream，Client为普通RPC请求，客户端发送一次普通的RPC请求，服务端通过流式响应多次发送数据集。

![7e983068967044febc79c9fa6e4fbaa3.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_decfd5fe44f84afc82341ac05ef4c4da.png "7e983068967044febc79c9fa6e4fbaa3.png")

### 1）server

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
/*
1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
    // 具体返回多少个response根据业务逻辑调整
    for n := 0; n <= 6; n++ {
        // 通过 send 方法不断推送数据
        err := stream.Send(&pb.StreamResponse{
            Pt: &pb.StreamPoint{
                Name:  r.Pt.Name,
                Value: r.Pt.Value + int32(n),
            },
        })
        if err != nil {
            return err
        }
        time.Sleep(time.Second)
    }
    // 返回nil表示已经完成响应
    return nil
}
```

在 Server，主要留意 stream.Send 方法。它看上去能发送 N 次？有没有大小限制？

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
type StreamService_ListServer interface {
    Send(*StreamResponse) error
    grpc.ServerStream
}

func (x *streamServiceListServer) Send(m *StreamResponse) error {
    return x.ServerStream.SendMsg(m)
}
```

通过阅读源码，可得知是 protoc 在生成时，根据定义生成了各式各样符合标准的接口方法。最终再统一调度内部的 SendMsg 方法，该方法涉及以下过程:

* •
  消息体（对象）序列化
* •
  压缩序列化后的消息体
* •
  对正在传输的消息体增加 5 个字节的 header
* •
  判断压缩+序列化后的消息体总字节长度是否大于预设的 maxSendMessageSize（预设值为 math.MaxInt32），若超出则提示错误
* •
  写入给流的数据集

### 2）client

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
/*
1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    // 调用获取stream
    stream, err := client.List(context.Background(), r)
    if err != nil {
        return err
    }
    // for循环获取服务端推送的消息
    for {
        // 通过 Recv() 不断获取服务端send()推送的消息
        resp, err := stream.Recv()
        // err==io.EOF则表示服务端关闭stream了
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
    }
    return nil
}
```

在 Client，主要留意 stream.Recv() 方法。什么情况下 io.EOF ？什么情况下存在错误信息呢?

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
type StreamService_ListClient interface {
    Recv() (*StreamResponse, error)
    grpc.ClientStream
}

func (x *streamServiceListClient) Recv() (*StreamResponse, error) {
    m := new(StreamResponse)
    if err := x.ClientStream.RecvMsg(m); err != nil {
        return nil, err
    }
    return m, nil
}
```

RecvMsg 会从流中读取完整的 gRPC 消息体，另外通过阅读源码可得知：

（1）RecvMsg 是阻塞等待的

（2）RecvMsg 当流成功/结束（调用了 Close）时，会返回 io.EOF

（3）RecvMsg 当流出现任何错误时，流会被中止，错误信息会包含 RPC 错误码。而在 RecvMsg 中可能出现如下错误：

* •
  io.EOF
* •
  io.ErrUnexpectedEOF
* •
  transport.ConnectionError
* •
  google.golang.org/grpc/codes

### 3）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:55149

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:55158
2022/11/03 09:35:03 resp: pj.name: gRPC Stream Client: List, pt.value: 1234
2022/11/03 09:35:04 resp: pj.name: gRPC Stream Client: List, pt.value: 1235
2022/11/03 09:35:05 resp: pj.name: gRPC Stream Client: List, pt.value: 1236
2022/11/03 09:35:06 resp: pj.name: gRPC Stream Client: List, pt.value: 1237
2022/11/03 09:35:07 resp: pj.name: gRPC Stream Client: List, pt.value: 1238
2022/11/03 09:35:08 resp: pj.name: gRPC Stream Client: List, pt.value: 1239
2022/11/03 09:35:09 resp: pj.name: gRPC Stream Client: List, pt.value: 1240
```

服务器流式 RPC 类似于一元 RPC，除了服务器返回消息流以响应客户端的请求。发送所有消息后，服务器的状态详细信息（状态代码和可选状态消息）和可选尾随元数据将发送到客户端。这样就完成了服务器端的处理。客户端在拥有服务器的所有消息后完成。

## 7、Client-side streaming RPC：客户端流式 RPC

客户端通过流式多次发送RPC请求给服务端，服务端发送一次响应给客户端。

![7e54a349f491420789281fbc54f362af.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_03487742906c4d9797c30b2a3b942669.png "7e54a349f491420789281fbc54f362af.png")

### 1）server

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
/*
1. for循环中通过stream.Recv()不断接收client传来的数据
2. err == io.EOF表示客户端已经发送完毕关闭连接了,此时在等待服务端处理完并返回消息
3. stream.SendAndClose() 发送消息并关闭连接(虽然在客户端流里服务器这边并不需要关闭 但是方法还是叫的这个名字，内部也只会调用Send())
*/
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
    // for循环接收客户端发送的消息
    for {
        // 通过 Recv() 不断获取客户端 send()推送的消息
        r, err := stream.Recv()
        // err == io.EOF表示已经获取全部数据
        if err == io.EOF {
            // SendAndClose 返回并关闭连接
            // 在客户端发送完毕后服务端即可返回响应
            return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
        }
        if err != nil {
            return err
        }
        log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
        time.Sleep(time.Second)
    }
    return nil
}
```

stream.SendAndClose：我们对每一个 Recv 都进行了处理，当发现 io.EOF (流关闭) 后，需要将最终的响应结果发送给客户端，同时关闭正在另外一侧等待的 Recv

### 2）client

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
/*
1. 建立连接并获取client
2. 获取 stream 并通过 Send 方法不断推送数据到服务端
3. 发送完成后通过stream.CloseAndRecv() 关闭stream并接收服务端返回结果
*/
func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    // 获取 stream
    stream, err := client.Record(context.Background())
    if err != nil {
        return err
    }

    for i := 0; i <= 6; i++ {
        // 通过 Send 方法不断推送数据到服务端
        err := stream.Send(r)
        if err != nil {
            return err
        }
    }

    // 发送完成后通过stream.CloseAndRecv() 关闭stream并接收服务端返回结果
    // (服务端则根据err==io.EOF来判断client是否关闭stream)
    resp, err := stream.CloseAndRecv()
    if err != nil {
        return err
    }
    log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
    return nil
}
```

stream.CloseAndRecv 和 stream.SendAndClose 是配套使用的流方法

### 3）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:57789
2022/11/03 11:59:31 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999
2022/11/03 11:59:32 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999
2022/11/03 11:59:33 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999
2022/11/03 11:59:34 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999
2022/11/03 11:59:35 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999
2022/11/03 11:59:36 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999
2022/11/03 11:59:37 stream.Recv pt.name: gRPC Stream Client: Record, pt.value: 9999

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:57793
2022/11/03 11:59:38 resp: pj.name: gRPC Stream Server: Record, pt.value: 1
```

## 8、Bidirectional streaming RPC：双向流式 RPC

双向流，由客户端发起流式的RPC方法请求，服务端以同样的流式RPC方法响应请求 首个请求一定是client发起，具体交互方法（谁先谁后，一次发多少，响应多少，什么时候关闭）根据程序编写方式来确定（可以结合协程）。

![c65a995764f745c6aeead21cf9b5f839.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_3d72309757634159928bb5ed26b27bef.png "c65a995764f745c6aeead21cf9b5f839.png")

### 1）server

一般是使用两个 Goroutine，一个接收数据，一个推送数据。最后通过 return nil 表示已经完成响应。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
/*
// 1. 建立连接 获取client
// 2. 通过client调用方法获取stream
// 3. 开两个goroutine（使用 chan 传递数据） 分别用于Recv()和Send()
// 3.1 一直Recv()到err==io.EOF(即客户端关闭stream)
// 3.2 Send()则自己控制什么时候Close 服务端stream没有close 只要跳出循环就算close了。 具体见https://github.com/grpc/grpc-go/issues/444
*/
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
    var (
        wg    sync.WaitGroup //任务编排
        msgCh = make(chan *pb.StreamPoint)
    )
    wg.Add(1)
    go func() {
        n := 0
        defer wg.Done()
        for v := range msgCh {
            err := stream.Send(&pb.StreamResponse{
                Pt: &pb.StreamPoint{
                    Name:  v.GetName(),
                    Value: int32(n),
                },
            })
            if err != nil {
                fmt.Println("Send error :", err)
                continue
            }
            n++
        }
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            r, err := stream.Recv()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Fatalf("recv error :%v", err)
            }
            log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
            msgCh <- &pb.StreamPoint{
                Name: "gRPC Stream Server: Route",
            }
        }
        close(msgCh)
    }()

    wg.Wait() //等待任务结束

    return nil
}
```

### 2）client

和服务端类似，不过客户端推送结束后需要主动调用 stream.CloseSend() 函数来关闭Stream。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
/*
1. 建立连接 获取client
2. 通过client获取stream
3. 开两个goroutine 分别用于Recv()和Send()
    3.1 一直Recv()到err==io.EOF(即服务端关闭stream)
    3.2 Send()则由自己控制
4. 发送完毕调用 stream.CloseSend()关闭stream 必须调用关闭 否则Server会一直尝试接收数据 一直报错...
*/
func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    var wg sync.WaitGroup
    // 调用方法获取stream
    stream, err := client.Route(context.Background())
    if err != nil {
        return err
    }

    // 开两个goroutine 分别用于Recv()和Send()
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            resp, err := stream.Recv()
            if err == io.EOF {
                fmt.Println("Server Closed")
                break
            }
            if err != nil {
                continue
            }
            log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
        }
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()

        for n := 0; n <= 6; n++ {
            err := stream.Send(r)
            if err != nil {
                log.Printf("send error:%v\n", err)
            }
            time.Sleep(time.Second)
        }

        // 发送完毕关闭stream
        err = stream.CloseSend()
        if err != nil {
            log.Printf("Send error:%v\n", err)
            return
        }
    }()

    wg.Wait()
    return nil
}
```

### 3）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:55108
2022/11/03 12:29:35 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/03 12:29:36 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/03 12:29:37 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/03 12:29:38 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/03 12:29:39 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/03 12:29:40 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111
2022/11/03 12:29:41 stream.Recv pt.name: gRPC Stream Client: Route, pt.value: 1111

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:55113
2022/11/03 12:29:35 resp: pj.name: gRPC Stream Server: Route, pt.value: 0
2022/11/03 12:29:36 resp: pj.name: gRPC Stream Server: Route, pt.value: 1
2022/11/03 12:29:37 resp: pj.name: gRPC Stream Server: Route, pt.value: 2
2022/11/03 12:29:38 resp: pj.name: gRPC Stream Server: Route, pt.value: 3
2022/11/03 12:29:39 resp: pj.name: gRPC Stream Server: Route, pt.value: 4
2022/11/03 12:29:40 resp: pj.name: gRPC Stream Server: Route, pt.value: 5
2022/11/03 12:29:41 resp: pj.name: gRPC Stream Server: Route, pt.value: 6
Server Closed
```

# 小结

客户端或者服务端都有对应的 推送或者 接收对象，我们只要 不断循环 Recv()或者 Send() 就能接收或者推送了！

gRPC Stream 和 goroutine 配合简直完美。通过 Stream 我们可以更加灵活的实现自己的业务。如 订阅，大数据传输等。

Client发送完成后需要手动调用Close()或者CloseSend()方法关闭stream，Server端则return nil就会自动 Close。

1）**ServerStream**

* •
  服务端处理完成后return nil代表响应完成
* •
  客户端通过 err == io.EOF判断服务端是否响应完成

2）**ClientStream**

* •
  客户端发送完毕通过CloseAndRecv关闭stream 并接收服务端响应
* •
  服务端通过 err == io.EOF判断客户端是否发送完毕，完毕后使用SendAndClose关闭 stream并返回响应。

3）**BidirectionalStream**

* •
  客户端服务端都通过stream向对方推送数据
* •
  客户端推送完成后通过CloseSend关闭流，通过err == io.EOF判断服务端是否响应完成
* •
  服务端通过err == io.EOF判断客户端是否响应完成,通过return nil表示已经完成响应

通过err == io.EOF来判定是否把对方推送的数据全部获取到了。

客户端通过CloseAndRecv或者CloseSend关闭 Stream，服务端则通过SendAndClose或者直接 return nil来返回响应。
