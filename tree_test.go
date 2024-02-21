package nanogo

import "testing"
import "github.com/stretchr/testify/assert"

func TestTreeNode_Put(t *testing.T) {
	root := &treeNode{name: "/", children: make([]*treeNode, 0)}
	root.Put("/usr/get/:id")
	root.Put("/usr/put/:id")
	root.Put("/usr/delete/:id")
	root.Put("/usr/patch/:id")
}
func TestTreeNode_Get(t *testing.T) {
	root := &treeNode{name: "/", children: make([]*treeNode, 0)}
	root.Put("/usr/get/:id")
	node := root.Get("/usr/get/1")
	t.Log(node)
	assert.Equal(t, ":id", node.name)
	root.Put("/usr/put/:name")
	node = root.Get("/usr/put/zarttic")
	assert.Equal(t, ":name", node.name)

}
