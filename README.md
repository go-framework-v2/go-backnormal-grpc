# go-backnormal

## 介绍

Go语言后端框架，结合go-zero等开源框架及工作使用的自研框架，定制更符合自身开发的go Web框架。

本项目主要介绍 go 的实施规范，并对规范进行了部分实现。另附 Go 学习指南一份：

[go 相对于 java 的优势](doc/go-advantage_lxb.md)

[java 人员学 go 的一般性问题](doc/go-common-qa_lxb.md)

官方：

- 交互式学习go: [A Tour of Go](https://go.dev/tour) [中文版](https://tour.go-zh.org)
- [Effective Go - The Go Programming Language (google.cn)](https://golang.google.cn/doc/effective_go)
- [用户手册：涵盖了语言，IDE，测试，数据库等方方面面](https://golang.google.cn/doc/)

民间：

- [Go 语言教程 | 菜鸟教程 (runoob.com)](https://www.runoob.com/go/go-tutorial.html)
- [Go语言入门教程，Golang入门教程（非常详细） (biancheng.net)](http://c.biancheng.net/golang/)

框架：

* [go-zero](https://go-zero.dev/docs/tasks/cli/mysql)

## go-backnormal项目目录结构

```
│ README.md        // 项目说明
├─doc              // 相关文档
└─src              // 源代码目录
 │ go.mod          // 模块依赖管理
 │ go.sum          // 模块依赖管理
 ├─api             // OpenAPI/Swagger 规范，JSON 模式文件，协议定义文件。
 ├─cmd             // 如有多个项目启动文件，可建此目录
 │ main.go         // 项目启动文件，如只有一个可直接挂到 src 下
 ├─config          // 运行所需的配置文件
 ├─internal        // 项目内部使用的代码
 │ ├─dao           // 访问数据库的代码
 │ ├─entity        // 数据实体
 │ └─service       // 服务实现
 ├─pkg             // 外部应用程序可以使用的库代码
 ├─test            // 额外的外部测试应用程序和测试数据
 ├─tool            // 支持工具。可以从 `/pkg` 和 `/internal` 使用。
 └─web             // 特定于 Web 组件、页面、模板等。
```

**注意**：虽然 src 目录不被外界建议，这里依然建议采用，原因如下：

- GOPAHT 已经成为过去式了。
- 可以保持项目根目录的清爽，在 git 服务器上很容易看到 Readme 的内容。
- go.mod 位与 src 目录内，会自动屏蔽 import 对 src 的引用。

参考[golang-standards / project-layout](https://gitcode.net/mirrors/golang-standards/project-layout/-/blob/master/README_zh.md)

## 关于公共包项目

因为 go 获取依赖项目是直接从源代码服务器上拉取的，所以 go.mod 文件中 `module` 指令必须与源码地址保持一致。

有下面的要点：

- 当主版本号变更时需为新的主版本号建立独立分支。
- 项目版本号为规范的格式：v[主版本号].[次版本号].[小版本号]
- 项目版本号以 源代码服务器上的 tag 形式体现.

`module` 示例：

```go-mod
// 0 或 1 版本 
module gitlab.cdel.local/[yourProject]
// 2 版本，下面的 v2 是一个新的分支
module gitlab.cdel.local/[yourProject]/v2
```
