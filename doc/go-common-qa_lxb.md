# java 人员学 go 的一般性问题

## 作用域

go 没有与 java 对等的 public, private, protected,friendly 关键字。要想跨文件使用 go 代码，名字的首字母必须大写。

## 包

一个目录下得所有文件的包名必须使用同一个包名，并保证包名和目录名相同。

引用本模块内部包时，需要以模块名作为前缀+包的完全路径。

引用本地其他模块时，可以在模块名前附加相对路径。

参考：[使用go module导入本地包 - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/109828249)

## 语言特性问题

### 如何处理继承

go 没有类，也没有继承的概念，go 用组合来模拟继承。

组合的方式可以把一个类的成员变量声明为另一个类类型，这样就可以使用另一个类的方法和属性，从而实现类似继承的行为。

示例： src/internal/service/inheritance

**要点**：注意匿名字段的使用，可以实现类似于继承的直接调用方式。

### 为已有类型扩展方法

**要点**：方法的接收者和方法必须位于同一个包。解决方法为已有类型定义一个新类型，注意是新类型而不是别名。

### 访问控制

示例： src/internal/service/extend
java: 需要使用 public, private 等关键字，go: 首字母大写便是 public, 否则便是private

### 接口

@@ -68,11 +56,9 @@ java 需要显式实现接口，而 go 则是隐式的，没有类似于 impleme

参考：[反射 (google.cn)](https://golang.google.cn/blog/laws-of-reflection)

### 指针

在使用指针时，需要注意：指针本身不为 nil 但指针的值可能为 nil，这经常会导致空指针 panic!。
**接口不要过大**

### 关闭资源
As [Rob Pike points out](https://go-proverbs.github.io/), "The bigger the interface, the weaker the abstraction."

### 对象比较

@@ -80,10 +66,44 @@ java 需要显式实现接口，而 go 则是隐式的，没有类似于 impleme

### 没有枚举

go 没有枚举定义，可用常量定义来模拟

### 深度复制

值类型的数据，默认全部都是深复制，Array、Int、String、Struct、Float，Bool。

引用类型的数据，默认全部都是浅复制，如指针，Slice，Map。

### 没有 synchronized

go 没有 synchronized， go 认为如果需要重入锁，那么代码是可以优化的。

参考：[Go sync.Mutex - 简书 (jianshu.com)](https://www.jianshu.com/p/9e5554617399)

### 如何处理继承

go 没有类，也没有继承的概念，go 用组合来模拟继承。

组合的方式可以把一个类的成员变量声明为另一个类类型，这样就可以使用另一个类的方法和属性，从而实现类似继承的行为。

示例： src/internal/service/inheritance

**要点**：注意匿名字段的使用，可以实现类似于继承的直接调用方式。

### 为已有类型扩展方法

**要点**：方法的接收者和方法必须位于同一个包。解决方法为已有类型定义一个新类型，注意是新类型而不是别名。

### 指针

在使用指针时，需要注意：指针本身不为 nil 但指针的值可能为 nil，这经常会导致空指针 panic!。

### 关闭资源

Java 使用 try-with-resources 语句来自动回收资源，go 使用 delay 来释放资源，参考：[Go语言defer（延迟执行语句） (biancheng.net)](http://c.biancheng.net/view/61.html)

示例：src/internal/service/free

### 错误与异常处理

相对于 java 的异常，go 会将问题分为两种情况进行处理：**意料之中的问题**和**意料之中的问题**。

## 框架问题

### 日志

go 内置的 log 库缺少级别和分割能力，本示例使用 [zap](https://github.com/uber-go/zap) + [lumberjack](https://github.com/natefinch/lumberjack)。[Frequently Asked Questions](https://github.com/uber-go/zap/blob/v1.24.0/FAQ.md), [go zap自定义日志输出格式](https://www.jianshu.com/p/fc90ea603ef2)

zap要点：

- 日志有两种输出方式：易用（可将对象直接输出为 json）的 zap.S() 和 高性能的 zap.L()

- 高性能的 logger 只支持结构化的日志输出

- 缺省输出会缓存，所以需要时常落盘： defer zap.L().Sync()

**时间格式化**：格式串必须是**go语言的诞生时间**，**01/02 03:04:05PM ‘06 -0700** ，我们常用的格式为： "2006-01-02 15:04:05.000"。 参考[Golang时间格式化 - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/145009400)

示例： src/tool/logger.go

**已知问题**：输出 json 时只输出 value 不输出 key。解决方法，自行序列化（性能不是很好）,因为 zap 之所以性能好，就是因为 It includes a reflection-free, zero-allocation JSON encoder.

### 配置、环境变量的读取

这里使用 [viper](https://github.com/spf13/viper) ，它提供了下面的功能（来源于 viper）：

- 设置默认值

- 从JSON，TOML，YAML，HCL，Envfile和Java属性properties配置文件读取

- 监控配置文件变更并重新加载（可选）

- 从环境变量读取

- 从远程配置系统（etcd 或 Consul）读取配置

- 从命令行标志读取

- 从缓冲区阅读

示例： src/config

### Web 服务

这里以 https://gin-gonic.com/ 为例进行演示

### 依赖注入

go 基本上不会有依赖注入问题 ，因为**Go 在设计上更倾向于明确的、显式的编程风格**

实际上是利用 go 的代码生成能力

### 注解

Java 支持注解， go 原生支持的不好，一般情况下不建议使用。理由是**Go 在设计上更倾向于明确的、显式的编程风格**。参考：[Go：我有注解，Java：不，你没有！ - 技术颜良 - 博客园 (cnblogs.com)](https://www.cnblogs.com/cheyunhua/p/15409847.html)

### json 处理

### 数据库编程

## 隐式接口

## 没有枚举

go 没有枚举定义，可用常量定义来模拟

## 深度复制

值类型的数据，默认全部都是深复制，Array、Int、String、Struct、Float，Bool。

引用类型的数据，默认全部都是浅复制，如指针，Slice，Map。

## json 处理

## 数据库编程

## 没有 synchronized

## 与 maven 对应的职能如何实现

- 依赖管理

- 打包静态资源