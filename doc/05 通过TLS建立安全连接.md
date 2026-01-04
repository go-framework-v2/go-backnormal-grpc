# 明文传输

先前的例子中 gRPC Client/Server 都是明文传输的，在明文通讯的情况下，你的请求就是裸奔的，有可能被第三方恶意篡改或者伪造为“非法”的数据。

我们抓个包查看一下：

![a3b019b7af2d4c529a3cd93bb5e999b4.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_ec0147f6e89e4738b45353fac18e49cd.png "a3b019b7af2d4c529a3cd93bb5e999b4.png")

![74984540f2fc4a97affec048fc11c8d1.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_f147ca2f64a8426dbacd6c11439bddb7.png "74984540f2fc4a97affec048fc11c8d1.png")

是明文传输，后面我们开始gRPC通过 TLS 证书建立安全连接，让数据能够加密处理，包括证书制作和CA签名校验等。

# TLS概述

传输层安全 (TLS) 对通过 Internet 发送的数据进行加密，以确保窃听者和黑客无法看到您传输的内容，这对于密码、信用卡号和个人通信等私人和敏感信息特别有用。

## 1、什么是TLS？

传输层安全 (TLS) 是一种 Internet 工程任务组 ( IETF ) 标准协议，可在两个通信计算机应用程序之间提供身份验证、隐私和数据完整性。它是当今使用最广泛部署的安全协议，最适合需要通过网络安全交换数据的 Web 浏览器和其他应用程序。这包括 Web 浏览会话、文件传输、虚拟专用网络 (VPN) 连接、远程桌面会话和 IP 语音 (VoIP)。最近，TLS 被集成到包括 5G 在内的现代蜂窝传输技术中，以保护整个无线电接入网络 ( RAN ) 的核心网络功能。

## 2、TLS的工作流程

TLS 使用客户端-服务器握手机制来建立加密和安全的连接，并确保通信的真实性。

* •
  通信设备交换加密功能。
* •
  使用数字证书进行身份验证过程以帮助证明服务器是它声称的实体。
* •
  发生会话密钥交换。在此过程中，客户端和服务器必须就密钥达成一致，以建立安全会话确实在客户端和服务器之间的事实——而不是在中间试图劫持会话的东西。

![4cacf692240c439fb8f5c12a19db4f37.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_7969894a4cf042cc991a855068580e0b.png "4cacf692240c439fb8f5c12a19db4f37.png")

# gRPC建立安全连接

## 1、概述

gRPC建立在HTTP/2协议之上，对TLS提供了很好的支持。当不需要证书认证时,可通过grpc.WithInsecure()选项跳过了对服务器证书的验证，没有启用证书的gRPC服务和客户端进行的是明文通信，信息面临被任何第三方监听的风险。为了保证gRPC通信不被第三方监听、篡改或伪造，可以对服务器启动TLS加密特性。

gRPC 内置了以下 encryption 机制：

* •
  SSL / TLS：通过证书进行数据加密；
* •
  ALTS：Google开发的一种双向身份验证和传输加密系统。

。只有运行在 Google Cloud Platform 才可用，一般不用考虑。

## 2、gRPC 加密类型

* •
  1）insecure connection：不使用TLS加密
* •
  2）server-side TLS：仅服务端TLS加密
* •
  3）mutual TLS：客户端、服务端都使用TLS加密

我们前面的例子都是明文传输的，使用的都是 insecure connection，通过指定 WithInsecure option 来建立 insecure connection，不建议在生产环境使用。

后面我们了解如何使用 TLS 来建立安全连接。

## 3、server-side TLS

### 1）流程

服务端 TLS 具体包含以下几个步骤：

* •
  制作证书，包含服务端证书和 CA 证书；
* •
  服务端启动时加载证书；
* •
  客户端连接时使用CA 证书校验服务端证书有效性。

也可以不使用 CA证书，即服务端证书自签名。

### 2）什么是CA？CA证书又是什么？

* •
  CA是Certificate Authority的缩写，也叫“证书授权中心”。它是负责管理和签发证书的第三方机构，作用是检查证书持有者身份的合法性，并签发证书，以防证书被伪造或篡改。

CA实际上是一个机构，负责“证件”印制核发。就像负责颁发身份证的公安局、负责发放行驶证、驾驶证的车管所。

* •
  CA 证书就是CA颁发的证书。我们常听到的数字证书就是CA证书,CA证书包含信息有:证书拥有者的身份信息，CA机构的签名，公钥和私钥。

。身份信息: 用于证明证书持有者的身份

。CA机构的签名: 用于保证身份的真实性

。公钥和私钥: 用于通信过程中加解密，从而保证通讯信息的安全性

### 3）什么是SAN？

SAN(Subject Alternative Name)是 SSL 标准 x509 中定义的一个扩展。使用了 SAN 字段的 SSL 证书，可以扩展此证书支持的域名，使得一个证书可以支持多个不同域名的解析。

我们在用go 1.15版本以上，用gRPC通过TLS建立安全连接时，会出现证书报错问题：

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
panic: rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: certificate
is not valid for any names, but wanted to match localhost"
```

造成这个panic的原因是从go 1.15 版本开始废弃 CommonName，我们没有使用官方推荐的 SAN 证书（默认是没有开启SAN扩展）而出现的错误，导致客户端和服务端无法建立连接。

### 4）目录结构

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
go-grpc-example
├── client
│   └──TLS_client
│   │   └──client.go
├── conf
│   └──ca.conf
│   └──server.conf
├── proto
│   └──search
│   │   └──search.proto
├── server
│   └──TLS_server
│   │   └──server.go
├── Makefile
```

### 5）生成CA根证书

在ca.conf里写入内容如下:

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

1.生成ca私钥，得到ca.key

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl genrsa -out ca.key 4096
```

openssl genrsa：生成RSA私钥，命令的最后一个参数，将指定生成密钥的位数，如果没有指定，默认512

2.生成ca证书签发请求，得到ca.csr

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
$ openssl req -new -sha256 -out ca.csr -key ca.key -config ca.conf
GB [CN]:
State or Province Name (full name) [ZheJiang]:
Locality Name (eg, city) [HuZhou]:
Organization Name (eg, company) [Step]:
linzyblog.netlify.app [linzyblog.netlify.app]:
```

这里一直回车就好了

openssl req：生成自签名证书，-new指生成证书请求、-sha256指使用sha256加密、-key指定私钥文件、-x509指输出证书、-days 3650为有效期，此后则输入证书拥有者信息

3.生成ca根证书，得到ca.crt

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl x509 -req -days 3650 -in ca.csr -signkey ca.key -out ca.crt
```

![d3e05b9a810e49bfa1c468bf8e9ffae0.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_77d47d2e03e44b0596d9527487b90bf2.png "d3e05b9a810e49bfa1c468bf8e9ffae0.png")

### 6）生成终端用户证书

在server.conf写入以下内容：

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
DNS.1   = go-grpc-example（这里很重要，客户端需要此字段做匹配）
IP      = 127.0.0.1
```

1.生成私钥，得到server.key

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl genrsa -out server.key 4096
```

2.生成证书签发请求，得到server.csr

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl req -new -sha256 -out server.csr -key server.key -config server.conf
```

这里也一直回车就好。

3.用CA证书生成终端用户证书，得到server.crt

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
```

![3a05c39cf5224b3d95114d31768f36da.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_ab44a3fbdf3c40609edcf9a3a7b3d0fb.png "3a05c39cf5224b3d95114d31768f36da.png")

### 7）server

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
const PORT = "8888"

func main() {
    // 根据服务端输入的证书文件和密钥构造 TLS 凭证
    c, err := credentials.NewServerTLSFromFile("./conf/server.pem", "./conf/server.key")
    if err != nil {
        log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
    }
    // 返回一个 ServerOption，用于设置服务器连接的凭据。
    // 用于 grpc.NewServer(opt ...ServerOption) 为 gRPC Server 设置连接选项
    lis, err := net.Listen("tcp", ":"+PORT) //创建 Listen，监听 TCP 端口
    if err != nil {
        log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
    }
    search.RegisterSearchServiceServer(s, &service{})

    s.Serve(lis)
}
```

### 8）client

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
const PORT = "8888"

func main() {
    // 根据客户端输入的证书文件和密钥构造 TLS 凭证。
    // 第二个参数 serverNameOverride 为服务名称。
    c, err := credentials.NewClientTLSFromFile("./conf/server.pem", "go-grpc-example")
    if err != nil {
        log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
    }
    // 返回一个配置连接的 DialOption 选项。
    // 用于 grpc.Dial(target string, opts ...DialOption) 设置连接选项
    conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }
    defer conn.Close()
    client := pb.NewSearchServiceClient(conn)
    resp, err := client.Search(context.Background(), &pb.SearchRequest{
        Request: "gRPC",
    })
    if err != nil {
        log.Fatalf("client.Search err: %v", err)
    }

    log.Printf("resp: %s", resp.GetResponse())
}
```

### 8）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:53981

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:54328
2022/11/03 19:35:10 resp: gRPC Server
```

抓个包再看看

![a41f9c64363c49b08a55a4aacffdf0bc.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_f4cca5a6e0a947189ee93e955ab10cfd.png "a41f9c64363c49b08a55a4aacffdf0bc.png")

## 4、mutual TLS

### 1）生成服务端证书

新增server.conf

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
DNS.1   = go-grpc-example（这里很重要，客户端需要此字段做匹配）
IP      = 127.0.0.1
```

/

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// 1. 生成私钥，得到server.key
openssl genrsa -out server.key 2048

//2. 生成证书签发请求，得到server.csr
openssl req -new -sha256 -out server.csr -key server.key -config server.conf

//3. 用CA证书生成终端用户证书，得到server.crt
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 365 -in server.csr -out se
rver.crt -extensions req_ext -extfile server.conf
```

### 2）生成客户端证书

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
// 1. 生成私钥，得到client.key
openssl genrsa -out client.key 2048

//2. 生成证书签发请求，得到client.csr
openssl req -new -key client.key -out client.csr 


//3. 用CA证书生成客户端证书，得到client.crt
 openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 365  -in client.csr -out client.crt
```

![87ad59abb90f448eb9ea003ef601bff3.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_e58c099653ab47cd9e51298c066c5fd2.png "87ad59abb90f448eb9ea003ef601bff3.png")

### 3）整理目录

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
conf
├── ca.conf
├── ca.crt
├── ca.csr
├── ca.key
├── client
│   ├── client.csr
│   ├── client.key
│   └── client.pem
├── server
│   ├── server.conf
|   └── server.crt
│   ├── server.csr
|   ├── server.key
└─server_side_TLS
```

### 4）server

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
const PORT = "8888"

func main() {
    // 公钥中读取和解析公钥/私钥对
    cert, err := tls.LoadX509KeyPair("./conf/server/server.crt", "./conf/server/server.key")
    if err != nil {
        fmt.Println("LoadX509KeyPair error", err)
        return
    }
    // 创建一组根证书
    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("./conf/ca.crt")
    if err != nil {
        fmt.Println("read ca pem error ", err)
        return
    }
    // 解析证书
    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        fmt.Println("AppendCertsFromPEM error ")
        return
    }

    c := credentials.NewTLS(&tls.Config{
        //设置证书链，允许包含一个或多个
        Certificates: []tls.Certificate{cert},
        //要求必须校验客户端的证书
        ClientAuth: tls.RequireAndVerifyClientCert,
        //设置根证书的集合，校验方式使用ClientAuth设定的模式
        ClientCAs: certPool,
    })
    s := grpc.NewServer(grpc.Creds(c))
    lis, err := net.Listen("tcp", ":"+PORT) //创建 Listen，监听 TCP 端口
    if err != nil {
        log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
    }
    //将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。
    //这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
    search.RegisterSearchServiceServer(s, &service{})

    //gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
    s.Serve(lis)
}
```

### 5）client

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
const PORT = "8888"

func main() {
    // 公钥中读取和解析公钥/私钥对
    cert, err := tls.LoadX509KeyPair("./conf/client/client.crt", "./conf/client/client.key")
    if err != nil {
        fmt.Println("LoadX509KeyPair error ", err)
        return
    }
    // 创建一组根证书
    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("./conf/ca.crt")
    if err != nil {
        fmt.Println("ReadFile ca.crt error ", err)
        return
    }
    // 解析证书
    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        fmt.Println("certPool.AppendCertsFromPEM error ")
        return
    }

    c := credentials.NewTLS(&tls.Config{
        Certificates: []tls.Certificate{cert},
        ServerName:   "go-grpc-example",
        RootCAs:      certPool,
    })

    conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }
    defer conn.Close()

    client := pb.NewSearchServiceClient(conn)
    resp, err := client.Search(context.Background(), &pb.SearchRequest{
        Request: "gRPC",
    })
    if err != nil {
        log.Fatalf("client.Search err: %v", err)
    }

    log.Printf("resp: %s", resp.GetResponse())
}
```

### 6）启动 & 请求

--javascripttypescriptshellbashsqljsonhtmlcssccppjavarubypythongorustmarkdown

```
# 启动服务端
$ go run server.go
API server listening at: 127.0.0.1:56036

# 启动客户端
$ go run client.go 
API server listening at: 127.0.0.1:56364
2022/11/03 20:21:55 resp: gRPC Server

# 更改ServerName为linzy
$ go run client.go 
API server listening at: 127.0.0.1:56424
2022/11/03 20:23:17 client.Search err: rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: cer
tificate is valid for go-grpc-example, not linzy"
```

抓个包看看

![88b6c6da1d9e40078ebadf5c047c1a42.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_4025a66e64b9410ea4dbe7dc00231b59.png "88b6c6da1d9e40078ebadf5c047c1a42.png")

![b1c173e848654a3c869e9f9d07ac0960.png](https://ucc.alicdn.com/pic/developer-ecology/maj75agy3asvu_d8064b39b2ed4760a2e9be0cc98657de.png "b1c173e848654a3c869e9f9d07ac0960.png")
