# Gee

Gee 是一个用 Go 实现的轻量级 Web 框架练习项目。
它提供了路由、中间件、动态路由参数、模板渲染、静态资源和错误恢复等基础能力。

## 功能

- `GET` / `POST` 路由注册
- 分组路由与分组中间件
- 动态路由参数：`:name`、`*filepath`
- 请求上下文封装：`Query`、`PostForm`、`Param`、`String`、`JSON`、`HTML`、`Data`
- 日志中间件
- `panic` 恢复中间件
- HTML 模板渲染
- 静态文件服务

## 结构

当前仓库是一个独立的 Go 模块，核心代码都在根目录：

- `gee.go`：引擎、路由分组、静态资源、模板等入口
- `router.go`：路由注册和路由匹配
- `trie.go`：前缀树实现
- `context.go`：请求上下文和响应封装
- `logger.go`：日志中间件
- `recovery.go`：异常恢复中间件

## 快速开始

确保本地已安装 Go，然后在项目根目录运行：

```bash
go run .
```

默认监听地址由示例程序决定，通常是 `:9999`。

## 示例

```go
package main

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.Default()

	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Gee\n")
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"name": c.Param("name"),
		})
	})

	r.Run(":9999")
}
```

## 常用能力

### 路由

- `/hello/:name` 可以匹配 `/hello/geektutu`
- `/assets/*filepath` 可以匹配 `/assets/css/style.css`

### 中间件

使用 `Use()` 注册中间件，执行顺序与分组顺序一致。`gee.Default()` 默认启用 `Recovery()`。

### 静态文件

可以通过 `Static()` 暴露本地目录，例如：

```go
r.Static("/assets", "./assets")
```

### 模板

先通过 `SetFuncMap()` 注册模板函数，再用 `LoadHTMLGlob()` 加载模板文件：

```go
r.SetFuncMap(template.FuncMap{
	"upper": strings.ToUpper,
})
r.LoadHTMLGlob("templates/*")
```

## 说明

这个仓库适合按源码顺序学习一个 Web 框架的核心实现。
如果你想从基础到完整版本浏览，可以直接顺着当前目录里的源码看实现细节。
