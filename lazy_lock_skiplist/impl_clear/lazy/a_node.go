package lazyskiplist

import "sync"

type Node struct {
	key, value  interface{}
	next        []*Node
	prev        *Node
	marked      bool
	fullyLinked bool
	lock        sync.Mutex
}

func newNode(key, value interface{}, level int) *Node {
	return &Node{
		key:         key,
		value:       value,
		next:        make([]*Node, level),
		marked:      false,
		fullyLinked: false}
}

func (node *Node) getLevel() int {
	return len(node.next)
}
