package main

import (
	"sync"
)

type SimpleMap struct {
	sync.RWMutex
	m map[int]bool
}

func NewSimpleMap() *SimpleMap {
	return &SimpleMap{
		m: make(map[int]bool, 0),
	}
}

func (h *SimpleMap) Get(id int) bool {
	h.RLock()
		v := h.m[id]
	h.RUnlock()
	return v
}

func (h *SimpleMap) Set(id int) {
	h.Lock()
		h.m[id] = true
	h.Unlock()
}

func (h *SimpleMap) Done(id int) {
	h.Lock()
		delete(h.m, id)
	h.Unlock()
}
