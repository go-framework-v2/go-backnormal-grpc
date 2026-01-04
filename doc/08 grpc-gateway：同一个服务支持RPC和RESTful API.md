# 一、gRPC-Gateway概述

## 1、简述

> 取自官方概述：
> grpc-gateway is a plugin of protoc. It reads gRPC service definition, and generates a reverse-proxy server which translates a RESTful JSON API into gRPC. This server is generated according to custom options in your gRPC definition.

gRPC-Gateway 是 `protoc 的插件`。它读取gRPC服务定义并`生成反向代理服务器，将 RESTful JSON API 转换为 gRPC`。该服务器是根据服务定义中的 `google.api.http` 注释生成的。

## 2、出现

etcd v3 改用 gRPC 后为了兼容原来的 API，同时要提供 HTTP/JSON 方式的API，为了满足这个需求，要么开发两套 API，要么实现一种转换机制，所以`grpc-gateway`诞生了。

* •
  通过protobuf的自定义option实现了一个`网关`，服务端同时开启gRPC和HTTP服务。
* •
  HTTP服务接收客户端请求后转换为grpc请求数据，获取响应后转为json数据返回给客户端。
* •
  当 HTTP 请求到达 gRPC-Gateway 时，它将 JSON 数据解析为 Protobuf 消息。`使用解析的 Protobuf 消息发出正常的 Go gRPC 客户端请求。`
* •
  Go gRPC 客户端将 Protobuf 结构编码为 `Protobuf 二进制格式`，然后将其发送到 gRPC 服务器。
* •
  gRPC 服务器处理请求并以 Protobuf 二进制格式返回响应。
* •
  Go gRPC 客户端将其解析为 Protobuf 消息，并将其返回到 gRPC-Gateway，后者将 Protobuf 消息编码为 JSON 并将其返回给原始客户端。

架构如下

![在这里插入图片描述](https://ucc.alicdn.com/images/user-upload-01/08317a546fc34ff5b7ed4eb2be0693f8.png?x-oss-process=image/resize,w_1400/format,webp "在这里插入图片描述")

# 二、准备工作

由于本实践偏向 Grpc+Grpc Gateway的方面，我们的需求是同一个服务端支持Rpc和Restful Api，那么就意味着`http2、TLS`等等的应用，功能方面就是一个服务端能够接受来自grpc和Restful Api的请求并响应。

本文示例代码已经上传到github：[点击跳转](https://github.com/Gopherlinzy/grpc-gateway-example)

## 1、目录结构

新建grpc-gateway-example文件夹，我们项目的初始目录目录如下：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
grpc-gateway-example/
├── certs
├── client
├── cmd
├── pkg
├── proto
│   ├── google
│   │   └── api
│   │   │   └── annotations.proto
│   │   │   └── http.proto
│   │   └── protobuf
│   │   │   └── descriptor.proto
├── server
└── Makefile
```

* •
  certs：存放证书凭证
* •
  client：客户端
* •
  cmd：存放 cobra 命令模块
* •
  pkg：第三方公共模块
* •
  proto：protobuf的一些相关文件（含.proto、pb.go、.pb.gw.go)，google/api中用于存放`annotations.proto、http.proto`、google/protobuf中用于存放`descriptor.proto`

  * •
    如果你生成Go代码的时候出现`File not found`，那一定就是找不到下面的文件。
  * •
    `annotations.proto和http.proto`文件需要手动从 [https://github.com/googleapis/googleapis/tree/master/google/api](https://github.com/googleapis/googleapis/tree/master/google/api)地址复制到自己的项目中或者手动复制代码！！
  * •
    `descriptor.proto`文件可以直接复制代码！！
* •
  server：服务端
* •
  Makefile：用于存放编译的代码。

## 2、环境准备

### 1）Protobuf

详细的请移步到[《gRPC（二）入门：Protobuf入门》](https://blog.csdn.net/weixin_46618592/article/details/127613897?spm=1001.2014.3001.5502)

### 2）gRPC

详细的请移步到[《gRPC（三）基础：gRPC快速入门》](https://blog.csdn.net/weixin_46618592/article/details/127630520?spm=1001.2014.3001.5502)

### 3）gRPC-Gateway

gRPC-Gateway 只是一个插件，只需要安装一下就可以了。这里建议科学上网：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

## 3、编写 IDL

### 1）google.api

`proto` 目录中有 `google/api` 目录，它用到了 google 官方提供的两个 api 描述文件，主要是针对 grpc-gateway 的 http 转换提供支持，定义了 Protocol Buffer 所扩展的 HTTP Option。

### 2）hello.proto

编写Demo的 `.proto `文件，我们在 proto目录下新建 hello.proto 文件，写入文件内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
syntax = "proto3";

package proto;
option go_package = "./proto/helloworld;helloworld";

import "proto/google/api/annotations.proto";

// 定义Hello服务
service Hello {
  // 定义SayHello方法
  rpc SayHello(HelloRequest) returns (HelloResponse) {
    // http option 网关
    option (google.api.http) = {
      post: "/hello_world"
      body: "*"
    };
  }
}

// HelloRequest 请求结构
message HelloRequest {
  string referer = 1;
}

// HelloResponse 响应结构
message HelloResponse {
  string message = 1;
}
```

在 hello.proto 文件中，引用了 `google/api/annotations.proto `，达到支持HTTP Option的效果

* •
  定义了一个 serviceRPC 服务 HelloWorld，在其内部定义了一个 `HTTP Option` 的POST方法，HTTP 响应路径为`/hello_world`。
* •
  定义message类型`HelloWorldRequest、HelloWorldResponse`，用于响应请求和返回结果。

> 每个方法都必须添加 `google.api.http` 注解后 gRPC-Gateway 才能生成对应 http 方法。
> 其中post为 HTTP Method，即 POST 方法， `/hello_world` 则是请求路径。

### 3）编译proto

在Makefile文件内输入以下内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
protoc:
        protoc  --go_out=.  --go-grpc_out=. --grpc-gateway_out=. ./proto/*.proto
```

* •
  Go Plugins 用于生成 .pb.go 文件
* •
  gRPC Plugins 用于生成 \_grpc.pb.go
* •
  gRPC-Gateway 则是 pb.gw.go

使用 `make protoc` 编译proto：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
➜ make protoc
protoc  --go_out=.  --go-grpc_out=. --grpc-gateway_out=. ./proto/*.proto
```

![在这里插入图片描述](https://ucc.alicdn.com/images/user-upload-01/a54704be7a9248fdb324b811dfda2a6f.png?x-oss-process=image/resize,w_1400/format,webp "在这里插入图片描述")

## 4、制作证书

详细的请移步到[《gRPC（五）进阶：通过TLS建立安全连接》](https://blog.csdn.net/weixin_46618592/article/details/127673304?spm=1001.2014.3001.5502)

在服务端支持Rpc和Restful Api，需要用到TLS，因此我们要先制作证书

进入certs目录，生成TLS所需的公钥密钥文件

### 1）生成CA根证书

在 `ca.conf` 文件并写入内容如下:

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
[ req ]
default_bits       = 4096
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = GB
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = ZheJiang
localityName                = Locality Name (eg, city)
localityName_default        = HuZhou
organizationName            = Organization Name (eg, company)
organizationName_default    = Step
commonName                  = linzyblog.netlify.app
commonName_max              = 64
commonName_default          = linzyblog.netlify.app
```

1. 生成ca私钥，得到ca.key

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl genrsa -out ca.key 4096
```

1. 生成ca证书签发请求，得到ca.csr

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
$ openssl req -new -sha256 -out ca.csr -key ca.key -config ca.conf
GB [CN]:
State or Province Name (full name) [ZheJiang]:
Locality Name (eg, city) [HuZhou]:
Organization Name (eg, company) [Step]:
linzyblog.netlify.app [linzyblog.netlify.app]:
```

1. 生成ca根证书，得到ca.crt

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl x509 -req -days 3650 -in ca.csr -signkey ca.key -out ca.crt
```

![在这里插入图片描述](https://ucc.alicdn.com/images/user-upload-01/2ccd0e1bdf80472eb461910d214e2f3f.png?x-oss-process=image/resize,w_1400/format,webp "在这里插入图片描述")

### 2）生成终端用户证书

在 `server.conf` 写入以下内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = ZheJiang
localityName                = Locality Name (eg, city)
localityName_default        = HuZhou
organizationName            = Organization Name (eg, company)
organizationName_default    = Step
commonName                  = CommonName (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = linzyblog.netlify.app

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1   = grpc-gateway-example
IP      = 127.0.0.1
```

1. 生成私钥，得到server.key

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl genrsa -out server.key 2048
```

1. 生成证书签发请求，得到server.csr

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl req -new -sha256 -out server.csr -key server.key -config server.conf
```

这里也一直回车就好。

1. 用CA证书生成终端用户证书，得到server.crt

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
```

![在这里插入图片描述](https://ucc.alicdn.com/images/user-upload-01/6c6d97ddb38749c68a3badb768d86f6e.png?x-oss-process=image/resize,w_1400/format,webp "在这里插入图片描述")

> 这样我们需要的证书凭证就足够了

# 三、命令行模块 cmd

## 1、Cobra介绍

官方文档：[点击跳转](https://grpc-ecosystem.github.io/grpc-gateway/)

`Cobra 是一个用于创建强大的现代 CLI 应用程序的库`。它提供了一个简单的界面来创建强大的现代 CLI 界面，类似于 git 和 go 工具。

Cobra 提供：

* •
  简易的子命令行模式
* •
  完全兼容 `POSIX` 的命令行模式(包括短版和长版）
* •
  嵌套的子命令
* •
  全局、本地和级联`flags`
* •
  使用Cobra很容易的生成应用程序和命令，使用 `cobra create appname`和 `cobra add cmdname`
* •
  提供智能提示
* •
  自动生成commands和flags的帮助信息
* •
  自动生成详细的 `help` 信息，如 app -help。
* •
  自动识别帮助 `flag、 -h，--help`。
* •
  自动生成应用程序在 bash 下命令自动完成功能。
* •
  自动生成应用程序的 man 手册。
* •
  命令行别名。
* •
  自定义 `help` 和 `usage` 信息。
* •
  可选的与 `viper` 的紧密集成。

## 2、概念

Cobra 建立在命令（commands）、参数（arguments ）、选项（flags）的结构之上。

* •
  **commands**：命令代表行为,一般表示 action，即运行的二进制命令服务。同时可以拥有子命令（children commands）
* •
  **arguments**：参数代表命令行参数。
* •
  **flags**：选项代表对命令行为的改变，即命令行选项。二进制命令的配置参数，可对应配置文件。参数可分为全局参数和子命令参数。

最好的命令行程序在实际使用时，就应该像在读一段优美的语句，能够更加直观的知道如何与用户进行交互。

执行命令行程序应该遵循一般的格式：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
#appname command  arguments
docker pull alpine:latest

#appname command flag
docker ps -a

#appname command flag argument
git commit -m "linzy"
```

## 3、安装

使用 Cobra 很容易。首先，用于go get安装最新版本的库。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go get -u github.com/spf13/cobra@latest
```

## 4、编写 server

在编写 cmd 时需要先用 server 进行测试关联，因此这一步我们先写 server.go 用于测试

在 server 模块下 新建 `server.go` 文件，写入测试内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package server

import (
    "log"
)

var (
    ServerPort  string
    CertName    string
    CertPemPath string
    CertKeyPath string
)

func Serve() (err error) {
    log.Println(ServerPort)

    log.Println(CertName)

    log.Println(CertPemPath)

    log.Println(CertKeyPath)

    return nil
}
```

## 5、编写 cmd

在cmd模块下 新建 `root.go` 文件，写入内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// rootCmd表示在没有任何子命令的情况下的基本命令
var rootCmd = &cobra.Command{
    // Command的用法，Use是一个行用法消息
    Use: "grpc",
    // Short是help命令输出中显示的简短描述
    Short: "Run the gRPC hello-world server",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}
```

当前 cmd 目录下继续 新建 `server.go` 文件，写入内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package cmd

import (
    "github.com/spf13/cobra"
    "grpc-gateway-example/server"
    "log"
)

// 创建附加命令
// 本地标签：在本地分配一个标志，该标志仅适用于该特定命令。
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Run the gRPC hello-world server",
    // 运行:典型的实际工作功能。大多数命令只会实现这一点；
    // 另外还有PreRun、PreRunE、PostRun、PostRunE等等不同时期的运行命令，但比较少用，具体使用时再查看亦可
    Run: func(cmd *cobra.Command, args []string) {
        defer func() {
            if err := recover(); err != nil {
                log.Println("Recover error : %v", err)
            }
        }()

        server.Serve()
    },
}

// 在 init() 函数中定义flags和处理配置。
func init() {
    // 我们定义了一个flag，值存储在&server.ServerPort中，长命令为--port，短命令为-p，，默认值为50052。
    // 命令的描述为server port。这一种调用方式成为Local Flags 本地标签
    serverCmd.Flags().StringVarP(&server.ServerPort, "port", "p", "50052", "server port")
    serverCmd.Flags().StringVarP(&server.CertPemPath, "cert-pem", "", "./certs/server.pem", "cert pem path")
    serverCmd.Flags().StringVarP(&server.CertKeyPath, "cert-key", "", "./certs/server.key", "cert key path")
    serverCmd.Flags().StringVarP(&server.CertName, "cert-name", "", "grpc-gateway-example", "server's hostname")

    // AddCommand向这父命令（rootCmd）添加一个或多个命令
    rootCmd.AddCommand(serverCmd)
}
```

## 6、启动 & 请求

我们在 `grpc-gateway-example` 目录下，新建文件main.go，写入内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import "grpc-gateway-example/cmd"

func main() {
    cmd.Execute()
}
```

当前目录下执行·go run main.go server·，查看输出是否为（此时应为默认值）：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
$ go run main.go server
2022/11/10 12:08:22 50052
2022/11/10 12:08:22 grpc-gateway-example                     
2022/11/10 12:08:22 ./certs/server.pem                       
2022/11/10 12:08:22 ./certs/server.key                       
```

执行`go run main.go server --port=8000 --cert-pem=test-pem --cert-key=test-key --cert-name=test-name`，检验命令行参数是否正确：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
$ go run main.go server --port=8000 --cert-pem=test-pem --cert-key=test-key --cert-name=test-name
2022/11/10 12:24:53 8000
2022/11/10 12:24:54 test-name
2022/11/10 12:24:54 test-pem 
2022/11/10 12:24:54 test-key
```

到这都无误，我们的 `cmd` 模块编写就正确了，下面开始我们的重点

## 7、目录结构

完成以上操作之后我们的目录是这样的结构，看看是不是缺少了：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
grpc-gateway-example/
├── certs
├── client
├── cmd                        // 命令行模块
│   ├── root.go
│   └── server.go
├── pkg
├── proto
│   ├── google
│   │   └── api
│   │   │   ├── annotations.proto
│   │   │   └── http.proto
│   │   └── protobuf
│   │   │   └── descriptor.proto
│   ├── helloworld
│   │   ├── hello.pb.go        // proto编译后文件
│   │   ├── hello.pb.gw.go    // gateway编译后文件
│   │   └── hello_grpc.pb.go// proto编译后接口文件
│   ├── hello.proto
├── server                    // GRPC服务端
│   └── server.go
└── Makefile
```

# 四、服务端模块 server

## 1、编写 hello.proto

在server目录下新建文件 `hello.go` ，写入文件内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package server

import (
    "context"
    "grpc-gateway-example/proto/helloworld"
)

type helloService struct {
    helloworld.UnimplementedHelloServer
}

func NewHelloService() *helloService {
    return &helloService{}
}

// ctx context.Context用于接受上下文参数
// r *pb.HelloWorldRequest用于接受protobuf的Request参数
func (h helloService) SayHello(ctx context.Context, r *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
    return &helloworld.HelloResponse{
        Message: "hello grpc-gateway",
    }, nil
}
```

## 2、编写 grpc.go

在 pkg 下新建 util 目录，新建 `grpc.go` 文件，写入内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package util

import (
    "google.golang.org/grpc"
    "net/http"
    "strings"
)

// 将gRPC请求和HTTP请求分别调用不同的handler处理。
func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
    if otherHandler == nil {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            grpcServer.ServeHTTP(w, r)
        })
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
            grpcServer.ServeHTTP(w, r)
        } else {
            otherHandler.ServeHTTP(w, r)
        }
    })
}
```

> `GrpcHandlerFunc` 函数是用于判断请求是来源于 Rpc 客户端还是 Restful Api 的请求，根据不同的请求注册不同的 ServeHTTP 服务； `r.ProtoMajor == 2` 也代表着请求必须基于 `HTTP/2`。
> 简而言之函数将gRPC请求和HTTP请求分别调用不同的handler处理。

如果不需要 TLS 建立安全链接，则可以使用`h2c`：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
    return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
            grpcServer.ServeHTTP(w, r)
        } else {
            otherHandler.ServeHTTP(w, r)
        }
    }), &http2.Server{})
}
```

## 3、编写 tls.go

在pkg下的 util 目录下，新建 `tls.go` 文件，写入内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package util

import (
    "crypto/tls"
    "golang.org/x/net/http2"
    "io/ioutil"
    "log"
)

// 用于处理从证书凭证文件（PEM），最终获取tls.Config作为HTTP2的使用参数
func GetTLSConfig(certPemPath, certKeyPath string) *tls.Config {
    var certKeyPair *tls.Certificate
    cert, _ := ioutil.ReadFile(certPemPath)
    key, _ := ioutil.ReadFile(certKeyPath)

    // 从一对PEM编码的数据中解析公钥/私钥对。成功则返回公钥/私钥对
    pair, err := tls.X509KeyPair(cert, key)
    if err != nil {
        log.Println("TLS KeyPair err: %v\n", err)
    }

    certKeyPair = &pair

    return &tls.Config{
        // tls.Certificate：返回一个或多个证书，实质我们解析PEM调用的X509KeyPair的函数声明
        // 就是func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (Certificate, error)，返回值就是Certificate
        Certificates: []tls.Certificate{*certKeyPair},
        // http2.NextProtoTLS：NextProtoTLS是谈判期间的NPN/ALPN协议，用于HTTP/2的TLS设置
        NextProtos: []string{http2.NextProtoTLS},
    }
}
```

`GetTLSConfig 函数是用于获取TLS配置`，在内部，我们读取了 server.key 和 server.pem 这类证书凭证文件。经过一系列处理获取 `tls.Config` 作为 HTTP2 的使用参数。

## 4、重新编写核心文件 server/server.go

修改server目录下的server.go文件，该文件是我们服务里的核心文件，写入内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package server

import (
    "context"
    "crypto/tls"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "grpc-gateway-example/pkg/util"
    "grpc-gateway-example/proto/helloworld"
    "log"
    "net"
    "net/http"
)

var (
    ServerPort  string
    CertName    string
    CertPemPath string
    CertKeyPath string
    EndPoint    string
)

func Serve() (err error) {
    EndPoint = ":" + ServerPort
    // 用于监听本地的网络地址通知
    // 它的函数原型func Listen(network, address string) (Listener, error)
    conn, err := net.Listen("tcp", EndPoint)
    if err != nil {
        log.Printf("TCP Listen err:%v\n", err)
    }

    // 通过util.GetTLSConfig解析得到tls.Config，传达给http.Server服务的TLSConfig配置项使用
    tlsConfig := util.GetTLSConfig(CertPemPath, CertKeyPath)
    srv := createInternalServer(conn, tlsConfig)

    log.Printf("gRPC and https listen on: %s\n", ServerPort)

    // NewListener将会创建一个Listener
    // 它接受两个参数，第一个是来自内部Listener的监听器，第二个参数是tls.Config（必须包含至少一个证书）
    if err = srv.Serve(tls.NewListener(conn, tlsConfig)); err != nil {
        log.Printf("ListenAndServe: %v\n", err)
    }

    return err
}

// 将认证的中间件注册进去, 前面所获取的tlsConfig仅能给HTTP使用
func createInternalServer(conn net.Listener, tlsConfig *tls.Config) *http.Server {
    var opts []grpc.ServerOption

    // 输入证书文件和服务器的密钥文件构造TLS证书凭证
    creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)
    if err != nil {
        log.Printf("Failed to create server TLS credentials %v", err)
    }

    // grpc.Creds()其原型为func Creds(c credentials.TransportCredentials) ServerOption
    // 该函数返回 ServerOption，它为服务器连接设置凭据
    opts = append(opts, grpc.Creds(creds))

    // 创建了一个没有注册服务的grpc服务端
    grpcServer := grpc.NewServer(opts...)

    // 注册grpc服务
    helloworld.RegisterHelloServer(grpcServer, NewHelloService())

    // 创建 grpc-gateway 关联组件
    // context.Background()返回一个非空的空上下文。
    // 它没有被注销，没有值，没有过期时间。它通常由主函数、初始化和测试使用，并作为传入请求的顶级上下文
    ctx := context.Background()
    // 从客户端的输入证书文件构造TLS凭证
    dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertName)
    if err != nil {
        log.Printf("Failed to create client TLS credentials %v", err)
    }
    // grpc.WithTransportCredentials 配置一个连接级别的安全凭据(例：TLS、SSL)，返回值为type DialOption
    // grpc.DialOption DialOption选项配置我们如何设置连接（其内部具体由多个的DialOption组成，决定其设置连接的内容）
    dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

    // 创建HTTP NewServeMux及注册grpc-gateway逻辑
    // runtime.NewServeMux：返回一个新的ServeMux，它的内部映射是空的；
    // ServeMux是grpc-gateway的一个请求多路复用器。它将http请求与模式匹配，并调用相应的处理程序
    gwmux := runtime.NewServeMux()

    // RegisterHelloWorldHandlerFromEndpoint：注册HelloWorld服务的HTTP Handle到grpc端点
    if err := helloworld.RegisterHelloHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
        log.Printf("Failed to register gw server: %v\n", err)
    }

    // http服务
    // 分配并返回一个新的ServeMux
    mux := http.NewServeMux()
    // 为给定模式注册处理程序
    mux.Handle("/", gwmux)

    return &http.Server{
        Addr:      EndPoint,
        Handler:   util.GrpcHandlerFunc(grpcServer, mux),
        TLSConfig: tlsConfig,
    }
}
```

## 5、server流程刨析

### 1）启动监听

`net.Listen("tcp", EndPoint)` 函数用于监听本地网络地址的监听。其函数原型`Listen(ctx context.Context, network, address string) (Listener, error)`

参数：

* •
  network：必须是tcp, tcp4, tcp6, unix或unixpacket。
* •
  address：对于TCP网络，如果address参数中的host为空或未指定的IP地址，则会自动返回一个可用的端口或者IP地址。

net.Listen("tcp", EndPoint)函数返回值是`Listener`：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
type Listener interface {
    // 接受等待并将下一个连接返回给Listener
    Accept() (Conn, error)

    // 关闭Listener
    Close() error

    // 返回 Listener 的网络地址。
    Addr() Addr
}
```

`net.Listen` 会返回一个监听器的结构体，返回接下来的动作，让其执行下一步的操作，可用执行以下操作Accept、Close、Addr。

### 2）获取TLSConfig

通过调用 `util.GetTLSConfig` 函数解析得到 `tls.Config`，通过传达给 `createInternalServer`函数完成 http.Server 服务的 `TLSConfig` 配置项使用。

### 3）创建内部服务

程序采用HTTP2、HTTPS，需要支持TLS，在启动 `grpc.NewServer()` 前需要将`serverOptions`（服务器选项，类似于中间件，可用设置例如凭证、编解码器和保持存活参数等选项。），而前面所获取的 tlsConfig 仅能给HTTP使用，因此第一步我们要创建 grpc 的 TLS 认证凭证。

1. 创建 grpc 的 TLS 认证凭证

引用 `google.golang.org/grpc/credentials` 第三方包，`credentials` 包实现gRPC库支持的各种凭据，这些凭据封装了客户机与服务器进行身份验证所需的所有状态，并进行各种断言，例如，关于客户机的身份、角色或是否授权进行特定调用。

我们调用 `NewServerTLSFromFile` 它能够从服务器的输入证书文件和密钥文件构造TLS凭据。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func NewServerTLSFromFile(certFile, keyFile string) (TransportCredentials, error) {
    // LoadX509KeyPair从一对文件中读取并解析一个公私钥对。文件中必须包含PEM编码的数据。
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }
    // NewTLS使用tls.Config来构建基于TLS的TransportCredentials（传输凭证）
    return NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}}), nil
}
```

1. grpc ServerOption

grpc.Creds() 其原型为`func Creds(c credentials.TransportCredentials) ServerOption`，返回一个为服务器连接设置凭据的ServerOption。

1. 创建 grpc 服务端

grpc.NewServer() 创建一个没有注册服务的grpc服务端，可以配置 `ServerOption`

1. 注册grpc服务

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// 注册grpc服务
helloworld.RegisterHelloServer(grpcServer, NewHelloService())
```

1. 创建 `grpc-gateway` 关联组件

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// context.Background()返回一个非空的空上下文。
// 它没有被注销，没有值，没有过期时间。它通常由主函数、初始化和测试使用，并作为传入请求的顶级上下文
ctx := context.Background()
// 从客户端的输入证书文件构造TLS凭证
dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertName)
if err != nil {
    log.Printf("Failed to create client TLS credentials %v", err)
}
// grpc.WithTransportCredentials 配置一个连接级别的安全凭据(例：TLS、SSL)，返回值为type DialOption
// grpc.DialOption DialOption选项配置我们如何设置连接（其内部具体由多个的DialOption组成，决定其设置连接的内容）
dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
```

1. 创建HTTP NewServeMux及注册 `grpc-gateway` 逻辑

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// 创建HTTP NewServeMux及注册grpc-gateway逻辑
// runtime.NewServeMux：返回一个新的ServeMux，它的内部映射是空的；
// ServeMux是grpc-gateway的一个请求多路复用器。它将http请求与模式匹配，并调用相应的处理程序
gwmux := runtime.NewServeMux()

// RegisterHelloWorldHandlerFromEndpoint：注册HelloWorld服务的HTTP Handle到grpc端点
if err := helloworld.RegisterHelloHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
    log.Printf("Failed to register gw server: %v\n", err)
}

// http服务
// 分配并返回一个新的ServeMux
mux := http.NewServeMux()
// 为给定模式注册处理程序
mux.Handle("/", gwmux)
```

1. 注册具体服务

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// RegisterHelloWorldHandlerFromEndpoint：注册HelloWorld服务的HTTP Handle到grpc端点
if err := helloworld.RegisterHelloHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
    log.Printf("Failed to register gw server: %v\n", err)
}
```

* •
  ctx：上下文
* •
  gwmux：`grpc-gateway` 的请求多路复用器
* •
  EndPoint：服务网络地址
* •
  dopts：配置好的安全凭据

### 4）创建Listener

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// NewListener将会创建一个Listener
// 它接受两个参数，第一个是来自内部Listener的监听器，第二个参数是tls.Config（必须包含至少一个证书）
if err = srv.Serve(tls.NewListener(conn, tlsConfig)); err != nil {
    log.Printf("ListenAndServe: %v\n", err)
}
```

### 5）服务接受请求

我们调用 `srv.Serve(tls.NewListener(conn, tlsConfig))`它是http.Server的方法，并且需要一个Listener作为参数。

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
func (srv *Server) Serve(l net.Listener) error {
    ...
    defer l.Close()
    ...
    baseCtx := context.Background()
    ...
    ctx := context.WithValue(baseCtx, ServerContextKey, srv)
    for {
        rw, e := l.Accept()
        ...
        c := srv.newConn(rw)
        c.setState(c.rwc, StateNew, runHooks) // before Serve can return
        go c.serve(ctx)
    }
}
```

它创建了一个 context.Background() 上下文对象，并调用 Listener 的 Accept 方法开始接受请求，在获取到连接数据后使用 newConn 创建连接对象，在最后使用goroutine的方式处理连接请求，完成请求后自动关闭连接。

# 五、验证功能

## 1、编写 client

在目录client下，创建 main.go 文件，新增以下内容：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
package main

import (
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "grpc-gateway-example/proto/helloworld"
    "log"
)

func main() {
    creds, err := credentials.NewClientTLSFromFile("./certs/server.pem", "grpc-gateway-example")
    if err != nil {
        log.Println("Failed to create TLS credentials %v", err)
        return
    }
    conn, err := grpc.Dial(":50052", grpc.WithTransportCredentials(creds))
    defer conn.Close()

    if err != nil {
        log.Println(err)
    }

    c := helloworld.NewHelloClient(conn)
    ct := context.Background()
    body := &helloworld.HelloRequest{
        Referer: "Grpc",
    }

    r, err := c.SayHello(ct, body)
    if err != nil {
        log.Println(err)
    }

    log.Println(r)
}
```

## 2、启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run main.go server
2022/11/10 16:34:06 gRPC and https listen on: 50052

# 启动客户端
$ go run client/main.go
2022/11/10 16:34:43 message:"hello grpc-gateway"
```

执行测试Restful Api，用POST方式访问[https://localhost](https://localhost/):50052/hello\_world
![在这里插入图片描述](https://ucc.alicdn.com/images/user-upload-01/6dd3db9196034a4291ff1592f1b14ce2.png?x-oss-process=image/resize,w_1400/format,webp "在这里插入图片描述")

# 六、最终目录结构

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
grpc-gateway-example/
├── certs                    //证书凭证
│   ├── ca.conf
│   ├── ca.crt
│   ├── ca.csr
│   ├── ca.key
│   ├── server.conf
│   ├── server.csr
│   ├── server.key
│   └── server.pem
├── client                    // 客户端
│   └── main.go
├── cmd                        // 命令行模块
│   ├── root.go
│   └── server.go
├── pkg                        // 第三方公共模块
│   └── util
│   │   │   ├── grpc.go
│   │   │   └── tls.go
├── proto
│   ├── google
│   │   └── api
│   │   │   ├── annotations.proto
│   │   │   └── http.proto
│   │   └── protobuf
│   │   │   └── descriptor.proto
│   ├── helloworld
│   │   ├── hello.pb.go        // proto编译后文件
│   │   ├── hello.pb.gw.go    // gateway编译后文件
│   │   └── hello_grpc.pb.go// proto编译后接口文件
│   ├── hello.proto
├── server                    // GRPC服务端
│   └── server.go
├── main.go
└── Makefile
```
