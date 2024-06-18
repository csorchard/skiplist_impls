package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// MaxLevel is the maximum level for the skip list
	MaxLevel int = 16
	// Probability is the probability factor for level generation
	Probability float32 = 0.5
)

type Node struct {
	value int
	next  []*Node
}

type SkipList struct {
	head  *Node
	level int // Current level of the skip list
}

func NewNode(value, level int) *Node {
	return &Node{
		value: value,
		next:  make([]*Node, level),
	}
}

func NewSkipList() *SkipList {
	return &SkipList{
		head:  NewNode(-1, MaxLevel),
		level: 0,
	}
}

func (sl *SkipList) randomLevel() int {
	lvl := 1
	for rand.Float32() < Probability && lvl < MaxLevel {
		lvl++
	}
	return lvl
}

func (sl *SkipList) Insert(value int) {
	pred := make([]*Node, MaxLevel)
	current := sl.head

	// 1. Find the predecessors
	for i := sl.level - 1; i >= 0; i-- {
		// < value is important. it ensures we only get predecessors
		for current.next[i] != nil && current.next[i].value < value {
			current = current.next[i]
		}
		pred[i] = current
	}

	current = current.next[0]
	if current == nil || current.value != value {
		lvl := sl.randomLevel()

		// Corner Case: fill the new higher levels with value as head.
		if lvl > sl.level {
			for i := sl.level; i < lvl; i++ {
				pred[i] = sl.head
			}
			sl.level = lvl
		}

		// 3. Link the new node
		newNode := NewNode(value, lvl)
		for i := 0; i < lvl; i++ {
			// kind of like linked list remove/update element
			newNode.next[i] = pred[i].next[i]
			pred[i].next[i] = newNode
		}
	}
}

func (sl *SkipList) Delete(value int) {
	update := make([]*Node, MaxLevel)
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].value < value {
			current = current.next[i]
		}
		update[i] = current
	}

	current = current.next[0]
	if current != nil && current.value == value {
		for i := 0; i < sl.level; i++ {
			if update[i].next[i] != current {
				break
			}
			// linkedin list delete
			update[i].next[i] = current.next[i]
		}
		for sl.level > 1 && sl.head.next[sl.level-1] == nil {
			sl.level--
		}
	}
}

func (sl *SkipList) Search(value int) bool {
	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].value < value {
			current = current.next[i]
		}
	}
	current = current.next[0]
	return current != nil && current.value == value
}

func main() {
	rand.Seed(time.Now().UnixNano())
	sl := NewSkipList()
	sl.Insert(3)
	sl.Insert(6)
	sl.Insert(7)
	sl.Insert(9)
	sl.Insert(12)
	sl.Insert(19)
	sl.Insert(17)

	fmt.Println("Search for 6:", sl.Search(6))
	fmt.Println("Search for 15:", sl.Search(15))

	sl.Delete(6)
	fmt.Println("Search for 6 after deletion:", sl.Search(6))
}
