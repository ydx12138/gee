package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 当前节点存储的路由段，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true，代表是通配符匹配
}

// 插入路由
// pattern是要插入的路由，parts是路由段切片，height是当前的层数(初始是0，是`/`)
func (n *node) insert(pattern string, parts []string, height int) {
	//递归出口
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//当前层级要插入的路由段
	//插入时，isWild的值取决于有没有:name或*name这样的动态路由
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//递归调用
	child.insert(pattern, parts, height+1)
}

// 查询路由
// parts是要查询的路由段切片，height是当前查到的层数，返回路由的最后一个节点
func (n *node) search(parts []string, height int) *node {
	//如果是
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	//查看下一层满足条件的节点，递归循环。注意：是在前往目标节点前就已经判断节点是否满足条件了，而非进入结点之后再判断
	children := n.matchChildren(part)
	//
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 第一个匹配成功的节点，用于插入(返回下一层第一个匹配的节点)
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找(返回下一层所有匹配的节点)
func (n *node) matchChildren(part string) []*node {
	var nodes []*node
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
