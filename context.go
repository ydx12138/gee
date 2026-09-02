package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 提供响应JSON/HTML/String/Data的方法
// 提供访问Query和PostForm参数的方法
type H map[string]interface{}
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	//request-info
	Path   string
	Method string
	Params map[string]string //解析后的参数
	//response-info
	StatusCode int
	//middleware
	handlers []HandlerFunc
	index    int
	//engine pointer
	engine *Engine
}

// 响应HTML
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

// 获取动态参数
func (c *Context) Param(key string) string {
	return c.Params[key]
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	context := &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
	return context
}

// 调用下一个中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// 表单参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 查询参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 立即响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 响应字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 响应JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.Status(code)
	if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
		//这里不响应其它状态码，是因为json.NewEncoder(c.Writer).Encode(obj)可能已经返回了200，在这里最后可以返回一些数据，但响应码已经不可改变了,所以推荐直接panic
		//panic(err)
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 响应[]byte
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) Fail(code int, err string) {
	//让执行链结束
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
