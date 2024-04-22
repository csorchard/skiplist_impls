package main

import (
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

const (
	MaxLevel = 16
	P        = 0.5
	NodeSize = uint64(unsafe.Sizeof(node{}))
)

type node struct {
	key  uint64
	next [MaxLevel]*node
}

type SkipList struct {
	head     *node
	maxLevel int
	arena    *Arena
}

func NewSkiplist(arena *Arena) *SkipList {
	head := (*node)(arena.Allocate(uint64(NodeSize)))
	if head == nil {
		panic("Arena out of space")
	}
	return &SkipList{
		head:     head,
		maxLevel: 1,
		arena:    arena,
	}
}

func (sl *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < P && level < MaxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(key uint64) {
	update := make([]*node, MaxLevel)
	x := sl.head
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for x.next[i] != nil && x.next[i].key < key {
			x = x.next[i]
		}
		update[i] = x
	}

	level := sl.randomLevel()
	if level > sl.maxLevel {
		for i := sl.maxLevel; i < level; i++ {
			update[i] = sl.head
		}
		sl.maxLevel = level
	}

	x = (*node)(sl.arena.Allocate(NodeSize))
	if x == nil {
		panic("Out of space")
	}
	x.key = key
	for i := 0; i < level; i++ {
		x.next[i] = update[i].next[i]
		update[i].next[i] = x
	}
}

func (sl *SkipList) Contains(key uint64) bool {
	x := sl.head
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for x.next[i] != nil && x.next[i].key < key {
			x = x.next[i]
		}
	}
	x = x.next[0]
	return x != nil && x.key == key
}

func main() {
	rand.Seed(time.Now().UnixNano())
	arena := NewArena(1024 * 1024) // 1MB arena
	sl := NewSkiplist(arena)

	// Example Usage
	sl.Insert(3)
	sl.Insert(6)
	sl.Insert(7)
	sl.Insert(9)
	sl.Insert(12)
	sl.Insert(19)
	sl.Insert(17)

	fmt.Println("Contains 6:", sl.Contains(6))
	fmt.Println("Contains 15:", sl.Contains(15))
}
