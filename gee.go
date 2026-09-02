package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

//ServeHTTP -> router.handle -> router.getRoute -> r.hanlders[key](c)

//GET -> addRoute -> router.addRoute -> parsePattern -> r.roots[method].insert -> r.hanlders[key] = handler

// gee.go是世界中心，不断从这里产出新方法，然后提取到其他文件里，只留下一个同名方法进行引用，可称之为中央集权制
type HandlerFunc func(c *Context)
type Engine struct {
	*RouterGroup  //engine作为整个分组的根
	router        *router
	groups        []*RouterGroup // store all groups
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string        // 前缀
	middlewares []HandlerFunc // 组上的中间件
	parent      *RouterGroup  // 父组
	engine      *Engine       // 组需要有能访问router的能力，所以有一个指向engine的指针
}

func Default() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} //engine.RouterGroup作为分组的根，在初始化的时候就放进去
	engine.Use(Recovery())
	return engine
}

// 初始化engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} //engine.RouterGroup作为分组的根，在初始化的时候就放进去
	return engine
}
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 加载全局资源，使其可以在响应中使用
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))

}

// 静态资源处理
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		//资源目录
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 加载静态资源，使其可以让用户通过路由直接访问
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// 添加中间件
func (group *RouterGroup) Use(midderware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, midderware...)
}

// 创建路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: []HandlerFunc{},
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
func (group *RouterGroup) GET(pattern string, handle HandlerFunc) {
	group.addRoute("GET", pattern, handle)
}
func (group *RouterGroup) POST(pattern string, handle HandlerFunc) {
	group.addRoute("POST", pattern, handle)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// 所有的请求会经过这里然后进行路由分配。这个方法的参数不能变，因为其是上层接口的实现
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//找出这个路由上的中间件，赋值给context
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}
