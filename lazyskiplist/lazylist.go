package lazyskiplist

import (
	"sync"
	"sync/atomic"
)

type LazyList[K any, V any] struct {
	// sentinels
	head, tail *node[K, V]
	less       func(a, b K) bool
	eq         func(a, b K) bool
}

type node[K any, V any] struct {
	key    K
	val    V
	next   atomic.Pointer[node[K, V]]
	marked atomic.Bool
	lock   sync.Mutex
}
