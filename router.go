package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //每种请求方式一个树
	hanlders map[string]HandlerFunc //
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		hanlders: make(map[string]HandlerFunc),
	}
}

// 输入路由，返回路由切片
// 一个路由里只允许有一个*，且在*之后的路由段是无效的
// 假设传入的路由是/hello/*action/login,返回的切片只会有["hello","*action"]
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	//如果树还不存在，就先初始化
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) // 树记录路由
	r.hanlders[key] = handler                 //map记录路由对应的HandlerFunc
}

// 根据路由，找到对应的handlerfunc方法进行处理
func (r *router) handle(c *Context) {
	n, map1 := r.getRoute(c.Method, c.Path)
	//把处理函数也加入到中间件列表的最后，中间件列表的顺序就是：中间件1-中间件2-中间件3-handle
	if n != nil {
		c.Params = map1
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.hanlders[key])
	} else {
		//如果没有handle，就在中间件列表最后放一个404提示
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	//开始执行
	c.Next()
}

// 提取动态参数，并返回最后一个node，这个node的pattern不是空的
// 譬如/hello/:name,请求是/hello/Jack,则返回的map是：name:Jack
// 譬如/hello/*name,请求是/hello/Jack/login,则返回的map是：name:Jack/login
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}
