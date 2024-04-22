package main

import (
	"sync/atomic"
	"unsafe"
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

func (a *Arena) Allocate(size uint64) unsafe.Pointer {
	offset := atomic.AddUint64(&a.offset, size) - size
	if offset+size > uint64(len(a.buf)) {
		return nil
	}
	return unsafe.Pointer(&a.buf[offset])
}
