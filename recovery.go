package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 错误恢复中间件
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// 打印用于调试的堆栈跟踪
// 获取当前 Goroutine 的调用栈，跳过前 3 层调用者（即 trace 本身、trace 的调用者、以及调用者的调用者），从真正业务代码处开始记录。
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // 跳过前3个

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)   //获取对应函数
		file, line := fn.FileLine(pc) //获取文件名和行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
