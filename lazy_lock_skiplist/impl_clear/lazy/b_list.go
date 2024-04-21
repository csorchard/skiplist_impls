package lazyskiplist

import (
	"fmt"
	"math/rand"
	"skiplist/lazy_lock_skiplist/impl_clear/lib"
	"sync/atomic"
)

var (
	MAX_LEVEL int = 16
	BRANCH    int = 4
)

type SkipList struct {
	head       *Node
	tail       *Node
	comparator lib.Comparator
	maxLevel   int
	size       int64
}

type OnUpdate func(old interface{}) interface{}

func NewLazySkipList(comparator lib.Comparator) *SkipList {
	head := &Node{next: make([]*Node, MAX_LEVEL)}
	tail := &Node{next: make([]*Node, MAX_LEVEL), prev: head}
	for i := range head.next {
		head.next[i] = tail
	}
	return &SkipList{
		head:       head,
		tail:       tail,
		comparator: comparator,
		maxLevel:   1}
}

func (this *SkipList) Size() int64 {
	return this.size
}

// Choose the new node's level, branching with p (1 / BRANCH) probability, with no regards to N (size of list)
func (this *SkipList) randomLevel() int {
	level := 1
	for level < MAX_LEVEL && rand.Intn(BRANCH) == 0 {
		level++
	}
	return level
}

func (this *SkipList) Get(key interface{}) (value interface{}, found bool) {
	pred := this.head
	for lv := MAX_LEVEL - 1; lv >= 0; lv-- {
		curr := pred.next[lv]
		for curr != this.tail && this.comparator(key, curr.key) > 0 {
			pred = curr
			curr = pred.next[lv]
		}

		if curr != this.tail && this.comparator(key, curr.key) == 0 {
			return curr.value, true
		}
	}
	return nil, false
}

func (this *SkipList) findNode(key interface{}, preds []*Node, succs []*Node) int {
	lFound := -1
	pred := this.head
	for lv := MAX_LEVEL - 1; lv >= 0; lv-- {
		curr := pred.next[lv]
		for curr != this.tail && this.comparator(key, curr.key) > 0 {
			pred = curr
			curr = pred.next[lv]
		}
		if lFound == -1 && curr != this.tail && this.comparator(key, curr.key) == 0 {
			lFound = lv
		}
		preds[lv] = pred
		succs[lv] = curr
	}
	return lFound
}

func (this *SkipList) tryPut(key, value interface{}, level int, preds []*Node, succs []*Node) bool {
	valid := true
	var prevPred *Node
	for lv := 0; valid && lv < level; lv++ {
		pred := preds[lv]
		succ := succs[lv]
		if pred != prevPred {
			pred.lock.Lock()
			defer pred.lock.Unlock()
			prevPred = pred
		}
		valid = !pred.marked && !succ.marked && pred.next[lv] == succ
	}
	if !valid {
		return false
	}

	node := newNode(key, value, level)

	node.prev = preds[0]

	for lv := 0; lv < level; lv++ {
		node.next[lv] = succs[lv]
		preds[lv].next[lv] = node
	}

	succs[0].prev = node

	node.fullyLinked = true
	return true
}

func (this *SkipList) Put(
	key, value interface{},
	onUpdate OnUpdate,
) (old interface{}, newbie interface{}, replaced bool) {

	level := this.randomLevel()

	preds := make([]*Node, MAX_LEVEL)
	succs := make([]*Node, MAX_LEVEL)

	for {
		lFound := this.findNode(key, preds, succs)
		if lFound != -1 {
			nodeFound := succs[lFound]
			if !nodeFound.marked {
				for !nodeFound.fullyLinked {
				}
				old = nodeFound.value
				if onUpdate != nil {
					newbie = onUpdate(old)
				} else {
					newbie = value
				}
				nodeFound.value = newbie
				return old, newbie, true
			}
			continue
		}

		if this.tryPut(key, value, level, preds, succs) {
			break
		}
	}

	atomic.AddInt64(&this.size, 1)

	return nil, value, false
}

func (this *SkipList) tryRemove(nodeToDelete *Node, preds []*Node, succs []*Node) bool {
	valid := true
	var prevPred *Node
	level := nodeToDelete.getLevel()
	for lv := 0; valid && lv < level; lv++ {
		pred := preds[lv]
		succ := succs[lv]
		if pred != prevPred {
			pred.lock.Lock()
			defer pred.lock.Unlock()
			prevPred = pred
		}
		valid = !pred.marked && pred.next[lv] == succ
	}
	if !valid {
		return false
	}

	for lv := level - 1; lv >= 0; lv-- {
		preds[lv].next[lv] = nodeToDelete.next[lv]
	}
	nodeToDelete.next[0].prev = preds[0]

	return true
}

func (this *SkipList) Remove(key interface{}) (value interface{}, ok bool) {
	var nodeToDelete *Node = nil
	isMarked := false
	preds := make([]*Node, MAX_LEVEL)
	succs := make([]*Node, MAX_LEVEL)
	for {
		lFound := this.findNode(key, preds, succs)
		if isMarked || (lFound != -1 && okToDelete(succs[lFound], lFound)) {
			if !isMarked {
				nodeToDelete = succs[lFound]
				nodeToDelete.lock.Lock()
				if nodeToDelete.marked {
					// someone else will remove.
					nodeToDelete.lock.Unlock()
					return nil, false
				}
				nodeToDelete.marked = true
				isMarked = true
			}

			if this.tryRemove(nodeToDelete, preds, succs) {
				nodeToDelete.lock.Unlock()
				break
			}
			nodeToDelete.lock.Unlock()
		} else {
			return nil, false
		}
	}
	atomic.AddInt64(&this.size, -1)
	return nodeToDelete.value, true
}

func okToDelete(node *Node, lFound int) bool {
	return node.fullyLinked && node.getLevel()-1 == lFound && !node.marked
}

func (this *SkipList) Print() {
	fmt.Print("[h] ")
	n := 0
	for i := this.head.next[0]; n < 100 && i.next[0] != nil; i = i.next[0] {
		if cap(i.next) > 1 {
			fmt.Printf("> [%v(%d)]", i.key, cap(i.next))
		} else {
			fmt.Printf("> [%v]", i.key)
		}
		if i.marked {
			fmt.Print("*")
		}
		fmt.Print(" ")
		n++
	}
	fmt.Print("> [t]\n")
}
