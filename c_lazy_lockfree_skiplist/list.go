package main

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const MaxLevel = 16
const P = 0.5 // Probability of level increment

type MarkableReference struct {
	node   *Node
	marked bool
}

type Node struct {
	value    int
	key      int
	next     []*atomic.Pointer[MarkableReference]
	topLevel int
}

type LockFreeSkipList struct {
	head *Node
	tail *Node
}

func newNode(value int, height int) *Node {
	n := &Node{
		value:    value,
		key:      value, // Using value as key for simplification
		next:     make([]*atomic.Pointer[MarkableReference], height+1),
		topLevel: height,
	}
	for i := range n.next {
		n.next[i] = new(atomic.Pointer[MarkableReference])
		mr := &MarkableReference{node: nil, marked: false}
		n.next[i].Store(mr)
	}
	return n
}

func NewLockFreeSkipList() *LockFreeSkipList {
	head := newNode(-1, MaxLevel)
	tail := newNode(int(^uint(0)>>1), MaxLevel)
	for i := range head.next {
		mr := &MarkableReference{node: tail, marked: false}
		head.next[i].Store(mr)
	}
	return &LockFreeSkipList{
		head: head,
		tail: tail,
	}
}

func (list *LockFreeSkipList) randomLevel() int {
	level := 0
	for rand.Float64() < P && level < MaxLevel {
		level++
	}
	return level
}

func (list *LockFreeSkipList) Add(value int) bool {
	topLevel := list.randomLevel()
	preds := make([]*Node, MaxLevel+1)
	succs := make([]*MarkableReference, MaxLevel+1)
	for {
		found := list.find(value, preds, succs)
		if found {
			return false // Value already present
		}

		newNode := newNode(value, topLevel)
		for i := 0; i <= topLevel; i++ {
			ref := succs[i]
			newNode.next[i].Store(&MarkableReference{node: ref.node, marked: false})
		}

		pred := preds[0]
		succ := succs[0]
		if !pred.next[0].CompareAndSwap(succ, &MarkableReference{node: newNode, marked: false}) {
			continue
		}

		for i := 1; i <= topLevel; i++ {
			for {
				pred = preds[i]
				succ = succs[i]
				if pred.next[i].CompareAndSwap(succ, &MarkableReference{node: newNode, marked: false}) {
					break
				}
				list.find(value, preds, succs)
			}
		}
		return true
	}
}

func (list *LockFreeSkipList) find(value int, preds []*Node, succs []*MarkableReference) bool {
	var pred, curr *Node
	var succ *MarkableReference
retry:
	for {
		pred = list.head
		for level := MaxLevel; level >= 0; level-- {
			curr = pred.next[level].Load().node
			for curr != list.tail && (curr.key < value || curr.next[level].Load().marked) {
				pred = curr
				curr = pred.next[level].Load().node
			}
			preds[level] = pred
			succs[level] = pred.next[level].Load()
		}
		succ = curr.next[0].Load()
		if succ.marked {
			continue retry // A marked node was reached; retry
		}
		return curr.key == value
	}
}

func (list *LockFreeSkipList) Contains(value int) bool {
	var curr *Node = list.head
	for level := MaxLevel; level >= 0; level-- {
		curr = curr.next[level].Load().node
		for curr.key < value {
			curr = curr.next[level].Load().node
		}
	}
	return curr.key == value && !curr.next[0].Load().marked
}

func (list *LockFreeSkipList) Delete(value int) bool {
	preds := make([]*Node, MaxLevel+1)
	succs := make([]*MarkableReference, MaxLevel+1)

	var victim *Node
	var isMarked bool
	var topLevel int

	for {
		// Step 1: Find the node to delete, and its predecessors and successors.
		found := list.find(value, preds, succs)
		if !found {
			victim = nil
			return false // If the node is not found, return false.
		}
		victim = succs[0].node
		if !isMarked {
			// Check if the node is already logically removed.
			topLevel = victim.topLevel
			victimNext := victim.next[topLevel].Load()
			if victim != succs[topLevel].node || victimNext.marked {
				return false // Victim node has changed, or already marked.
			}

			// Mark the node logically from the top level down.
			isMarked = true
			for level := topLevel; level >= 0; level-- {
				succ := victim.next[level].Load()
				if !succ.marked {
					// Attempt to mark the node at each level.
					expected := &MarkableReference{node: succ.node, marked: false}
					newMark := &MarkableReference{node: succ.node, marked: true}
					if !victim.next[level].CompareAndSwap(expected, newMark) {
						isMarked = false
						break // If marking fails, break and retry
					}
				}
			}
			if !isMarked {
				continue // If marking failed, restart from finding the victim
			}
		}

		// Step 2: Physically remove the node.
		for level := topLevel; level >= 0; level-- {
			pred := preds[level]
			succ := succs[level]
			next := victim.next[level].Load().node
			// Attempt to unlink the victim node.
			if !pred.next[level].CompareAndSwap(&MarkableReference{node: succ.node, marked: false}, &MarkableReference{node: next, marked: false}) {
				isMarked = false // If CAS failed, set isMarked to false and retry
				break
			}
		}
		if !isMarked {
			continue // If any CAS failed, retry
		}
		return true
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	lfs := NewLockFreeSkipList()
	lfs.Add(10)
	lfs.Add(20)
	lfs.Add(30)

	added := lfs.Add(20) // should return false because 20 is already in the list
	println("Added 20 again:", added)

	contains := lfs.Contains(20)
	println("Contains 20:", contains)

	deleted := lfs.Delete(20)
	println("Deleted 20:", deleted)

	contains = lfs.Contains(20)
	println("Contains 20 after deletion:", contains)
}
