package main

import (
	"errors"
	"sync/atomic"
)

type Arena struct {
	offset uint64
	buf    []byte
}

func NewArena(size uint64) *Arena {
	return &Arena{
		buf: make([]byte, size),
	}
}

func (a *Arena) Alloc(size uint64) (uint64, error) {
	offset := atomic.AddUint64(&a.offset, size) - size
	if offset+size > uint64(len(a.buf)) {
		return 0, errors.New("arena: out of space")
	}
	return offset, nil
}
