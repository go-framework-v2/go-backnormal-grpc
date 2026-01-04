## 1. 使用Swagger/OpenAPI生成API文档

### 安装必要工具

```
# 安装swag工具
go install github.com/swaggo/swag/cmd/swag@latest

# 安装Gin的Swagger支持库
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 项目结构配置

```
/myproject
  /docs
    docs.go       # 自动生成
    swagger.json  # 自动生成
    swagger.yaml  # 自动生成
  /api
    handlers.go   # 你的API处理器
  main.go        # 主程序入口
  go.mod
```

## 2. 添加API注释

### 主程序注释 (main.go)

```
// @title Go微服务API文档
// @version 1.0
// @description 这是一个使用Gin框架构建的微服务API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API支持
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
    // 你的主程序代码
}
```

### 路由处理器注释 (api/handlers.go)

```
package api

import "github.com/gin-gonic/gin"

// GetUser godoc
// @Summary 获取用户信息
// @Description 根据ID获取用户详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
    // 处理器实现
}

// CreateUser godoc
// @Summary 创建新用户
// @Description 创建新用户
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "用户信息"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
    // 处理器实现
}
```

### 数据结构注释


```
// User 用户数据结构
type User struct {
    // 用户ID
    ID int `json:"id" example:"1"`
    // 用户名
    Name string `json:"name" example:"张三"`
    // 用户邮箱
    Email string `json:"email" example:"user@example.com"`
    // 创建时间
    CreatedAt time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
    // 错误码
    Code int `json:"code" example:"400"`
    // 错误信息
    Message string `json:"message" example:"无效的请求参数"`
}
```

## 3. 生成文档

```
# 在项目根目录执行
swag init -g main.go --parseDependency --parseInternal --parseDepth 2
```

## 4. 集成Swagger UI

```
import (
    _ "yourproject/docs" // 生成的docs包
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
)

func main() {
    r := gin.Default()
  
    // 添加Swagger路由
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
  
    // 其他路由配置...
  
    r.Run(":8080")
}
```

## 5. 访问API文档

启动服务后，访问以下URL查看API文档：

```
http://localhost:8080/swagger/index.html
```

## 6. 高级配置

### 分组API文档

```
// 在注释中使用@BasePath指定不同组的基础路径
// AdminGroup godoc
// @title 管理员API
// @version 1.0
// @description 管理员专用API
// @BasePath /admin
func AdminGroup(r *gin.RouterGroup) {
    // 管理员路由
}
```

### 安全定义

```
// 主程序注释中添加安全定义
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// 在需要认证的路由上添加注释
// @Security ApiKeyAuth
```

### 自定义Swagger UI配置

```
r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, 
    ginSwagger.DefaultModelsExpandDepth(-1),
    ginSwagger.PersistAuthorization(true),
    ginSwagger.DocExpansion("none"),
))
```

## 7. 生产环境注意事项

1. **禁用Swagger UI**：在生产环境中应禁用或保护Swagger UI

   ```
   if gin.Mode() != gin.ReleaseMode {
       r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
   }
   ```
2. **文档版本控制**：将生成的docs目录纳入版本控制
3. **CI/CD集成**：在构建流程中加入文档生成步骤
4. **文档导出**：可以将swagger.json导出到API管理平台

   ```
   curl http://localhost:8080/swagger/doc.json > swagger.json
   ```

这套API文档配置方案提供了完整的API描述能力，支持自动生成交互式文档，并能与代码保持同步更新，大大提高了API的可维护性和开发效率。
