package main

import (
	"math/rand"
	"sync"
	"time"
)

const MaxLevel = 16

type Node struct {
	value int
	next  []*Node
	lock  sync.Mutex

	marked      bool
	fullyLinked bool
	topLevel    int
}

type LazySkipList struct {
	head, tail *Node
}

func NewNode(item int, height int) *Node {
	node := &Node{
		value:    item,
		next:     make([]*Node, height+1),
		topLevel: height,
	}
	return node
}

func NewLazySkipList() *LazySkipList {
	head := NewNode(-1, MaxLevel)
	tail := NewNode(int(^uint(0)>>1), MaxLevel)
	for i := range head.next {
		head.next[i] = tail
	}
	return &LazySkipList{
		head: head,
		tail: tail,
	}
}

func (list *LazySkipList) randomLevel() int {
	lvl := 1
	for rand.Float32() < 0.5 && lvl < MaxLevel {
		lvl++
	}
	return lvl
}

func (list *LazySkipList) find(key int, preds []*Node, succs []*Node) int {
	lFound := -1
	pred := list.head
	for level := MaxLevel; level >= 0; level-- {
		curr := pred.next[level]
		for curr.value < key {
			pred = curr
			curr = pred.next[level]
		}
		if lFound == -1 && curr.value == key {
			lFound = level
		}
		preds[level] = pred
		succs[level] = curr
	}
	return lFound
}

func (list *LazySkipList) Add(item int) bool {
	topLevel := list.randomLevel()
	preds := make([]*Node, MaxLevel+1)
	succs := make([]*Node, MaxLevel+1)
	for {
		lFound := list.find(item, preds, succs)
		if lFound != -1 {
			nodeFound := succs[lFound]
			if !nodeFound.marked {
				for !nodeFound.fullyLinked {
				}
				return false
			}
			continue
		}

		highestLocked := -1
		valid := true
		for level := 0; valid && level <= topLevel; level++ {
			pred := preds[level]
			succ := succs[level]
			pred.lock.Lock()
			highestLocked = level
			valid = !pred.marked && !succ.marked && pred.next[level] == succ
		}
		if !valid {
			for level := 0; level <= highestLocked; level++ {
				preds[level].lock.Unlock()
			}
			continue
		}

		newNode := NewNode(item, topLevel)
		for level := 0; level <= topLevel; level++ {
			newNode.next[level] = succs[level]
			preds[level].next[level] = newNode
		}
		newNode.fullyLinked = true
		for level := 0; level <= highestLocked; level++ {
			preds[level].lock.Unlock()
		}
		return true
	}
}

func (list *LazySkipList) Remove(item int) bool {
	// Similar to Add, with deletion logic.
	// Implementation omitted for brevity.
	return true
}

func (list *LazySkipList) Contains(item int) bool {
	// Similar to find logic in Add.
	// Implementation omitted for brevity.
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())
	list := NewLazySkipList()
	success := list.Add(10)
	println("Added 10:", success)
	success = list.Add(20)
	println("Added 20:", success)
	found := list.Contains(10)
	println("Contains 10:", found)
	removed := list.Remove(10)
	println("Removed 10:", removed)
}
