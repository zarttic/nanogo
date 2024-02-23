package nanogo

import "strings"

type treeNode struct {
	name       string
	children   []*treeNode
	routerName string
	isEnd      bool
}

// Put  path: /usr/get/:id
// Put 方法将给定路径按照"/"分割后，逐层插入到树节点中。
// 如果遇到已存在的子节点，则移动到该子节点；否则创建新的子节点，并根据是否为路径的最后一部分标记isEnd属性。
func (t *treeNode) Put(path string) {
	// 初始化根节点引用
	root := t
	// 按"/"分割输入路径
	strs := strings.Split(path, "/")

	for index, name := range strs {
		// 跳过路径起始的"/"
		if index == 0 {
			continue
		}

		// 获取当前节点的所有子节点
		children := t.children
		// 标记是否存在匹配的子节点
		isMatch := false

		// 遍历子节点寻找匹配项
		for _, node := range children {
			if node.name == name {
				isMatch = true
				// 若找到匹配项，将当前节点移动至该子节点
				t = node
				break
			}
		}

		// 若未找到匹配项，则创建新节点并添加至当前节点的子节点列表中
		if !isMatch {
			isEnd := false
			// 判断当前层级是否为路径的最后一级
			if index == len(strs)-1 {
				isEnd = true
			}

			newNode := &treeNode{name: name, children: make([]*treeNode, 0), isEnd: isEnd}
			children = append(children, newNode)
			t.children = children
			// 移动当前节点至新创建的子节点
			t = newNode
		}
	}

	// 将当前节点重置为根节点
	t = root
}

// Get 根据路径获取树节点 path: /usr/get/1
func (t *treeNode) Get(path string) *treeNode {
	// 将路径按"/"分割成字符串数组
	strs := strings.Split(path, "/")
	// 初始化routerName为空字符串
	routerName := ""
	// 遍历字符串数组
	for index, name := range strs {
		// 跳过第一个元素
		if index == 0 {
			continue
		}
		// 获取当前节点的子节点
		children := t.children
		// 初始化isMatch为false
		isMatch := false
		// 遍历子节点
		for _, node := range children {
			// 如果子节点的名称与当前字符串匹配，或者子节点的名称为"*"，或者子节点的名称包含":"，则满足条件
			if node.name == name ||
				node.name == "*" ||
				strings.Contains(node.name, ":") {
				// 更新isMatch为true
				isMatch = true
				// 更新routerName为当前节点的名称
				routerName += "/" + node.name
				// 更新当前节点的routerName为routerName
				node.routerName = routerName
				// 更新当前节点为子节点
				t = node
				// 如果已经遍历到最后一个字符串，则返回当前节点
				if index == len(strs)-1 {
					return node
				}
				// 跳出内层循环
				break
			}
		}
		// 如果没有找到匹配的子节点，则继续判断是否存在通配符子节点
		if !isMatch {
			// 遍历子节点
			for _, node := range children {
				// 如果子节点的名称为"**"，则满足条件
				if node.name == "**" {
					// 更新routerName为当前节点的名称
					routerName += "/" + node.name
					// 更新当前节点的routerName为routerName
					node.routerName = routerName
					// 返回当前节点
					return node
				}
			}
		}
	}
	// 如果没有找到匹配的节点，则返回nil
	return nil
}
