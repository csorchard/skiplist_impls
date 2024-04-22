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
)

type Node struct {
	key  uint64
	next [MaxLevel]uint64 // Store offsets instead of pointers
}

type Skiplist struct {
	head     uint64 // Store offset of head
	maxLevel int
	arena    *Arena
}

func NewSkiplist(arena *Arena) *Skiplist {
	headOffset, err := arena.Alloc(uint64(unsafe.Sizeof(Node{})))
	if err != nil {
		panic("failed to allocate head node")
	}
	head := (*Node)(unsafe.Pointer(&arena.buf[headOffset]))
	*head = Node{} // Initialize the head node

	return &Skiplist{
		head:     headOffset,
		maxLevel: 1,
		arena:    arena,
	}
}

func (sl *Skiplist) getNode(offset uint64) *Node {
	return (*Node)(unsafe.Pointer(&sl.arena.buf[offset]))
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Float64() < P && level < MaxLevel {
		level++
	}
	return level
}

func (sl *Skiplist) Insert(key uint64) {
	update := [MaxLevel]uint64{}
	x := sl.head
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for sl.getNode(x).next[i] != 0 && sl.getNode(sl.getNode(x).next[i]).key < key {
			x = sl.getNode(x).next[i]
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

	nodeOffset, err := sl.arena.Alloc(uint64(unsafe.Sizeof(Node{})))
	if err != nil {
		panic("failed to allocate new node")
	}
	newNode := sl.getNode(nodeOffset)
	newNode.key = key

	for i := 0; i < level; i++ {
		newNode.next[i] = sl.getNode(update[i]).next[i]
		sl.getNode(update[i]).next[i] = nodeOffset
	}
}

func (sl *Skiplist) Contains(key uint64) bool {
	x := sl.head
	for i := sl.maxLevel - 1; i >= 0; i-- {
		for sl.getNode(x).next[i] != 0 && sl.getNode(sl.getNode(x).next[i]).key < key {
			x = sl.getNode(x).next[i]
		}
	}
	x = sl.getNode(x).next[0]
	return x != 0 && sl.getNode(x).key == key
}

func main() {
	rand.Seed(time.Now().UnixNano())
	arena := NewArena(1024 * 1024) // 1MB arena
	sl := NewSkiplist(arena)

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
