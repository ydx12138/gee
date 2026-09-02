# Gee

Gee 是一个用 Go 实现的小型 Web 框架练习项目。项目按 `day0` 到 `day7-1` 逐步演进，从最基础的 `net/http` 封装开始，逐步加入上下文、路由、动态参数、路由分组、中间件、模板渲染、静态资源服务和错误恢复。

这个仓库适合用来学习 Web 框架的核心结构：请求如何进入框架、如何匹配路由、如何把请求和响应封装到上下文里，以及中间件链如何串联起来。

## 功能特性

- HTTP 路由注册：支持 `GET`、`POST`
- 上下文封装：提供 `Query`、`PostForm`、`Param`、`String`、`JSON`、`HTML`、`Data` 等常用方法
- 动态路由：支持 `:name` 参数和 `*filepath` 通配符
- 路由分组：支持分组前缀和分组中间件
- 中间件机制：支持按顺序执行多个处理函数
- 日志中间件：记录请求处理信息
- Recovery 中间件：捕获 panic 并返回 500，避免服务崩溃
- HTML 模板：支持模板加载和自定义模板函数
- 静态资源：支持通过路由访问本地静态文件

## 目录结构

```text
.
|-- day0/       # 最基础的 net/http 示例
|-- day1-1/     # 框架雏形
|-- day1-2/     # 初步封装 Engine
|-- day2-1/     # 上下文 Context
|-- day3-1/     # 前缀树路由与动态参数
|-- day4-1/     # 路由分组
|-- day5-1/     # 中间件
|-- day6-1/     # HTML 模板与静态资源
`-- day7-1/     # 错误恢复与当前完整版本
```

每个 `day*` 目录都是一个相对独立的 Go 模块，里面包含示例入口 `main.go` 和对应阶段的 `gee` 框架代码。建议从 `day7-1` 查看当前功能最完整的版本。

## 快速开始

确保本地已经安装 Go，然后进入最新版本目录运行：

```bash
cd day7-1
go run .
```

服务默认监听：

```text
http://localhost:9999
```

访问首页：

```bash
curl http://localhost:9999/
```

预期输出：

```text
Hello Geektutu
```

也可以访问下面的地址测试 Recovery 中间件：

```bash
curl http://localhost:9999/panic
```

该接口会主动触发 panic，框架会捕获错误并返回 500 响应，而不是让服务直接退出。

## 使用示例

```go
package main

import (
	"gee"
	"net/http"
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

## 核心概念

### Engine

`Engine` 是框架入口，实现了 `http.Handler` 接口。所有请求都会进入 `ServeHTTP`，再由路由模块找到对应的处理函数。

### Context

`Context` 封装了一次 HTTP 请求和响应，简化了参数读取和响应写入。业务处理函数只需要操作 `*gee.Context`，不用直接处理底层的 `http.ResponseWriter` 和 `*http.Request`。

### Router

路由使用前缀树实现，支持普通路径、命名参数和通配符参数。例如：

- `/hello/:name` 可以匹配 `/hello/geektutu`
- `/assets/*filepath` 可以匹配 `/assets/css/style.css`

### Middleware

中间件和最终处理函数组成一条执行链，通过 `Context.Next()` 向后执行。`gee.Default()` 默认启用 Recovery 中间件。

## 开发说明

如果你想逐步理解框架实现，可以按目录顺序阅读：

```text
day0 -> day1-1 -> day1-2 -> day2-1 -> day3-1 -> day4-1 -> day5-1 -> day6-1 -> day7-1
```

每个阶段只引入少量新能力，适合对照提交历史和代码变化学习。

## 许可协议

当前仓库还没有显式声明开源许可证。正式开源前建议补充 `LICENSE` 文件，例如 MIT、Apache-2.0 或其他适合你的协议。
